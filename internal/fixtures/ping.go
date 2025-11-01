package fixtures

import (
	"net"
	"testing"
	"time"

	"github.com/bruli/pinger/internal/domain"
	"github.com/stretchr/testify/require"
)

type PingBuilder struct {
	ResourceName      *string
	Target            *net.IPAddr
	Interval, Timeout *time.Duration
	WarnMs            *float64
}

func (b PingBuilder) Build(t *testing.T) *domain.Ping {
	name := setData("name", b.ResourceName)
	target := setData(BuildTarget(), b.Target)
	interval := setData(time.Second, b.Interval)
	timeout := setData(time.Second, b.Timeout)
	warnMs := setData(100.0, b.WarnMs)

	ping, err := domain.NewPing(name, target, interval, timeout, warnMs)
	require.NoError(t, err)
	return ping
}

func BuildTarget() net.IPAddr {
	return net.IPAddr{IP: net.ParseIP("127.0.0.1")}
}
