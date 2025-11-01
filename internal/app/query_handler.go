package app

import (
	"context"
	"fmt"
)

type Query interface {
	Name() string
}

type QueryHandler interface {
	Handle(ctx context.Context, query Query) (any, error)
}

type InvalidQueryError struct {
	expected, had string
}

func (i InvalidQueryError) Error() string {
	return fmt.Sprintf("invalic query, expected: %s, had: %s", i.expected, i.had)
}

func NewInvalidQueryError(expected string, had string) *InvalidQueryError {
	return &InvalidQueryError{expected: expected, had: had}
}
