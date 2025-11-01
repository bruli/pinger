package domain_test

import (
	"net"
	"testing"
	"time"

	"github.com/bruli/pinger/internal/domain"
	"github.com/bruli/pinger/internal/fixtures"
	"github.com/stretchr/testify/require"
)

func TestNewPing(t *testing.T) {
	name := "resourceName"
	target := net.IPAddr{IP: net.ParseIP("127.0.0.1")}
	interval := time.Second
	timeout := time.Second
	warnMs := 100.0
	type args struct {
		name     string
		target   net.IPAddr
		interval time.Duration
		timeout  time.Duration
		warnMs   float64
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "with an invalid resourceName, then it returns an invalid ping resourceName error",
			args: args{
				name: "",
			},
			expectedErr: domain.ErrInvalidPingResourceName,
		},
		{
			name: "with an invalid target, then it returns an invalid ping target error",
			args: args{
				name:   name,
				target: net.IPAddr{},
			},
			expectedErr: domain.ErrInvalidPingTarget,
		},
		{
			name: "with an invalid interval, then it returns an invalid ping interval error",
			args: args{
				name:     name,
				target:   target,
				interval: 0,
			},
			expectedErr: domain.ErrInvalidPingInterval,
		},
		{
			name: "with an invalid timeout, then it returns an invalid ping timeout error",
			args: args{
				name:     name,
				target:   target,
				interval: interval,
				timeout:  0,
			},
			expectedErr: domain.ErrInvalidPingTimeout,
		},
		{
			name: "with an invalid warn ms, then it returns an invalid ping warn ms error",
			args: args{
				name:     name,
				target:   target,
				interval: interval,
				timeout:  timeout,
			},
			expectedErr: domain.ErrInvalidPingWanMs,
		},
		{
			name: "with valid data, then it returns a valid struct",
			args: args{
				name:     name,
				target:   target,
				interval: interval,
				timeout:  timeout,
				warnMs:   warnMs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(`Given a NewPing function,
		when is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewPing(tt.args.name, tt.args.target, tt.args.interval, tt.args.timeout, tt.args.warnMs)
			if err != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}
			require.Equal(t, tt.args.name, got.ResourceName())
			require.Equal(t, tt.args.target, got.Target())
			require.Equal(t, tt.args.interval, got.Interval())
			require.Equal(t, tt.args.timeout, got.Timeout())
			require.Equal(t, tt.args.warnMs, got.WarnMs())
			require.Len(t, got.Events(), 0)
		})
	}
}

func TestPing_AddLatency(t *testing.T) {
	type args struct {
		lat domain.Latency
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus string
	}{
		{
			name: "with a latency below warn ms, then it returns a ready status",
			args: args{
				lat: domain.Latency(100.0),
			},
			expectedStatus: domain.ReadyStatus.String(),
		},
		{
			name: "with a latency major warn ms, then it returns a degraded status",
			args: args{
				lat: domain.Latency(150.0),
			},
			expectedStatus: domain.DegradedStatus.String(),
		},
	}
	for _, tt := range tests {
		t.Run(`Given a Ping struct,
		when AddLatency method is called `+tt.name, func(t *testing.T) {
			t.Parallel()
			warm := 120.0
			p := fixtures.PingBuilder{WarnMs: &warm}.Build(t)
			p.AddLatency(tt.args.lat)
			events := p.Events()
			require.Len(t, events, 1)
			event := events[0]
			require.IsType(t, domain.PingEvent{}, event)
			require.Equal(t, tt.expectedStatus, event.(domain.PingEvent).Status)
			require.Equal(t, tt.args.lat.Float64(), event.(domain.PingEvent).Latency)
			require.Equal(t, domain.PingResultEventName, event.(domain.PingEvent).Name)
			require.False(t, event.CreatedAt().IsZero())
		})
	}
}

func TestPing_AddFail(t *testing.T) {
	t.Run(`Given a Ping struct,
	when AddFail method is called,
	then it returns an event with fail status`, func(t *testing.T) {
		p := fixtures.PingBuilder{}.Build(t)
		p.AddFail()
		events := p.Events()
		require.Len(t, events, 1)
		event := events[0]
		require.IsType(t, domain.PingEvent{}, event)
		require.Equal(t, domain.FailStatus.String(), event.(domain.PingEvent).Status)
		require.Equal(t, domain.PingResultEventName, event.(domain.PingEvent).Name)
		require.False(t, event.CreatedAt().IsZero())
	})
}
