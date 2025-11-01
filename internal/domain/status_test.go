package domain_test

import (
	"testing"

	"github.com/bruli/pinger/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "with an invalid status, then it returns an invalid status error",
			args: args{
				status: "invalid",
			},
			expectedErr: domain.ErrInvalidStatus,
		},
		{
			name: "with a ready status, then it returns a valid status",
			args: args{
				status: domain.ReadyStatus.String(),
			},
		},
		{
			name: "with a degraded status, then it returns a valid status",
			args: args{
				status: domain.DegradedStatus.String(),
			},
		},
		{
			name: "with a fail status, then it returns a valid status",
			args: args{
				status: domain.FailStatus.String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a ParseStatus function,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.ParseStatus(tt.args.status)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.Equal(t, tt.args.status, got.String())
		})
	}
}
