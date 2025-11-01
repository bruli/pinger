package app

import (
	"context"

	"github.com/bruli/pinger/internal/domain"
)

//go:generate go tool moq -out zmock_executor.go . PingExecutor
type PingExecutor interface {
	Execute(ctx context.Context, p *domain.Ping) (domain.Latency, error)
}
