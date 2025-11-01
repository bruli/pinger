package yaml

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bruli/pinger/internal/domain"
	"gopkg.in/yaml.v3"
)

type Ping struct {
	ResourceName string `yaml:"name"`
	Target,
	Interval, Timeout string
	WarnMs float64 `yaml:"warn_ms"`
}

type PingsRepository struct {
	data []byte
}

func (p PingsRepository) Find(ctx context.Context) ([]*domain.Ping, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var pings []Ping
		if err := yaml.Unmarshal(p.data, &pings); err != nil {
			return nil, err
		}
		return buildDomain(pings)
	}
}

func buildDomain(data []Ping) ([]*domain.Ping, error) {
	pings := make([]*domain.Ping, len(data))
	for i, d := range data {
		interval, err := time.ParseDuration(d.Interval)
		if err != nil {
			return nil, fmt.Errorf("invalid interval: %w", err)
		}
		timeout, err := time.ParseDuration(d.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout: %w", err)
		}
		ip, err := net.LookupIP(d.Target)
		if err != nil {
			return nil, err
		}
		if len(ip) == 0 {
			return nil, fmt.Errorf("no ip found for %s", d.Target)
		}
		d.Target = ip[0].String()
		p, err := domain.NewPing(d.ResourceName, net.IPAddr{IP: ip[0]}, interval, timeout, d.WarnMs)
		if err != nil {
			return nil, err
		}
		pings[i] = p
	}
	return pings, nil
}

func NewPingsRepository(data []byte) *PingsRepository {
	return &PingsRepository{data: data}
}
