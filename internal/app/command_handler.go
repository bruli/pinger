package app

import (
	"context"
	"fmt"

	"github.com/bruli/pinger/internal/domain"
)

type Command interface {
	Name() string
}

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) ([]domain.Event, error)
}

type InvalidCommandError struct {
	expected, had string
}

func (i InvalidCommandError) Error() string {
	return fmt.Sprintf("invalic command, expected: %s, had: %s", i.expected, i.had)
}

func NewInvalidCommandError(expected string, had string) *InvalidCommandError {
	return &InvalidCommandError{expected: expected, had: had}
}
