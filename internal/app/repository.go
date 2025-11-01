package app

import (
	"context"

	"github.com/bruli/pinger/internal/domain"
)

//go:generate go tool moq -out zmock_repository.go . PingRepository
type PingRepository interface {
	Find(ctx context.Context) ([]*domain.Ping, error)
}
