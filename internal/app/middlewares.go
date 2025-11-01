package app

import (
	"context"

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
			switch {
			case err != nil:
				log.Error().Err(err).Msg("error handling command")
			default:
				for _, e := range events {
					ev, _ := e.(domain.PingEvent)
					log.Info().Msgf("Executed. Target %q, latency %v, status %q", ev.AggregateRootID(), ev.Latency, ev.Status)
				}
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
