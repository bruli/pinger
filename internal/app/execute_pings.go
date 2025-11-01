package app

import (
	"context"

	"github.com/bruli/pinger/internal/domain"
)

const ExecutePingsCmdName = "executePings"

type ExecutePingsCommand struct {
	Ping *domain.Ping
}

func (e ExecutePingsCommand) Name() string {
	return ExecutePingsCmdName
}

type ExecutePings struct {
	exec PingExecutor
}

func (e ExecutePings) Handle(ctx context.Context, cmd Command) ([]domain.Event, error) {
	co, ok := cmd.(ExecutePingsCommand)
	if !ok {
		return nil, NewInvalidCommandError(ExecutePingsCmdName, cmd.Name())
	}
	ping := co.Ping
	lat, err := e.exec.Execute(ctx, ping)
	switch {
	case err != nil:
		ping.AddFail()
	default:
		ping.AddLatency(lat)
	}
	return ping.Events(), nil
}

func NewExecutePings(exec PingExecutor) *ExecutePings {
	return &ExecutePings{exec: exec}
}
