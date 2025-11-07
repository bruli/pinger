package app

import (
	"context"
	"errors"

	"github.com/bruli/pinger/internal/domain"
	"github.com/rs/zerolog"
)

type (
	CommandHandlerMiddleware func(h CommandHandler) CommandHandler
	QueryHandlerMiddleware   func(h QueryHandler) QueryHandler
)

type CommandHandlerFunc func(ctx context.Context, cmd Command) ([]domain.Event, error)

func (c CommandHandlerFunc) Handle(ctx context.Context, cmd Command) ([]domain.Event, error) {
	return c(ctx, cmd)
}

type QueryHandlerFunc func(ctx context.Context, query Query) (any, error)

func (c QueryHandlerFunc) Handle(ctx context.Context, query Query) (any, error) {
	return c(ctx, query)
}

func NewLogCommandHandlerMiddleware(log *zerolog.Logger) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]domain.Event, error) {
			events, err := h.Handle(ctx, cmd)
			if err != nil {
				log.Err(err).Msgf("failed to handle command %q", cmd.Name())
			}
			return events, err
		})
	}
}

func NewEventBusCommandHandlerMiddleware(bus EventBus, log *zerolog.Logger) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		return CommandHandlerFunc(func(ctx context.Context, cmd Command) ([]domain.Event, error) {
			events, err := h.Handle(ctx, cmd)

			errs := make([]error, 0)

			for _, event := range events {
				if err = bus.Dispatch(ctx, event); err != nil {
					errs = append(errs, err)
				}
			}
			if len(errs) > 0 {
				err = errors.Join(errs...)
				log.Err(err).Msg("failed to dispatch events")
			}

			return events, err
		})
	}
}

func NewLogQueryHandlerMiddleware(logger *zerolog.Logger) QueryHandlerMiddleware {
	return func(h QueryHandler) QueryHandler {
		return QueryHandlerFunc(func(ctx context.Context, query Query) (any, error) {
			result, err := h.Handle(ctx, query)
			if err != nil {
				logger.Error().Err(err).Msg("error handling query")
			}
			return result, err
		})
	}
}

func NewCommandHandlerMultiMiddleware(middlewares ...CommandHandlerMiddleware) CommandHandlerMiddleware {
	return func(h CommandHandler) CommandHandler {
		handler := h
		for _, m := range middlewares {
			handler = m(handler)
		}
		return CommandHandlerFunc(handler.Handle)
	}
}
