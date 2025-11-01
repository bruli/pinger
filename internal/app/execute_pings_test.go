package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bruli/pinger/internal/app"
	"github.com/bruli/pinger/internal/domain"
	"github.com/bruli/pinger/internal/fixtures"
	"github.com/stretchr/testify/require"
)

func TestExecutePings_Handle(t *testing.T) {
	ctx := context.Background()
	warm := float64(5)
	ping := fixtures.PingBuilder{WarnMs: &warm}.Build(t)
	type args struct {
		ctx context.Context
		cmd app.Command
	}
	tests := []struct {
		name                 string
		args                 args
		expectedErr, execErr error
		latency              domain.Latency
		expectedStatus       domain.Status
	}{
		{
			name: "with an invalid command, then it returns an invalid command error",
			args: args{
				ctx: ctx,
				cmd: InvalidCommand{},
			},
			expectedErr: app.InvalidCommandError{},
		},
		{
			name: "and executor returns an error, then it returns a failed event",
			args: args{
				ctx: ctx,
				cmd: app.ExecutePingsCommand{Ping: ping},
			},
			execErr:        errors.New("error"),
			expectedStatus: domain.FailStatus,
		},
		{
			name: "and executor returns a major latency, then it returns a degraded event",
			args: args{
				ctx: ctx,
				cmd: app.ExecutePingsCommand{Ping: ping},
			},
			latency:        domain.Latency(150),
			expectedStatus: domain.DegradedStatus,
		},
	}
	for _, tt := range tests {
		t.Run(`Given a ExecutePings command handler,
		when Handler method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			exec := &app.PingExecutorMock{}
			exec.ExecuteFunc = func(ctx context.Context, p *domain.Ping) (domain.Latency, error) {
				return tt.latency, tt.execErr
			}
			handler := app.NewExecutePings(exec)
			events, err := handler.Handle(tt.args.ctx, tt.args.cmd)
			if err != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
				return
			}
			require.Len(t, events, 1)
			ev, ok := events[0].(domain.PingEvent)
			require.True(t, ok)
			require.Equal(t, tt.expectedStatus.String(), ev.Status)
		})
	}
}

type InvalidCommand struct{}

func (i InvalidCommand) Name() string {
	return "invalid"
}
