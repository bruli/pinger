package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/bruli/pinger/internal/domain"
)

type UnknownEventToDispatchError struct {
	event string
}

func (u UnknownEventToDispatchError) Error() string {
	return fmt.Sprintf("event %q is not declared to dispatch", u.event)
}

type EventListener interface {
	Listen(ctx context.Context, ev domain.Event) error
}

type EventBus map[string][]EventListener

func NewEventBus() EventBus {
	return make(map[string][]EventListener)
}

func (e EventBus) Subscribe(ev domain.Event, listeners ...EventListener) {
	e[ev.EventName()] = listeners
}

func (e EventBus) Dispatch(ctx context.Context, ev domain.Event) error {
	errs := make([]error, 0)
	list, ok := e[ev.EventName()]
	if !ok {
		return UnknownEventToDispatchError{event: ev.EventName()}
	}
	for _, l := range list {
		if err := l.Listen(ctx, ev); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
