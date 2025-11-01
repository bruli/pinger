package domain

import (
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
)

const PingResultEventName = "ping.result"

var (
	ErrInvalidPingResourceName = errors.New("invalid ping resourceName")
	ErrInvalidPingTarget       = errors.New("invalid ping target")
	ErrInvalidPingTimeout      = errors.New("invalid ping timeout")
	ErrInvalidPingInterval     = errors.New("invalid ping interval")
	ErrInvalidPingWanMs        = errors.New("invalid ping warn ms")
)

type Latency float64

func (l Latency) Float64() float64 {
	return float64(l)
}

type Ping struct {
	BasicAggregateRoot
	resourceName      string
	target            net.IPAddr
	interval, timeout time.Duration
	warnMs            float64
}

func NewPing(resourceName string, target net.IPAddr, interval time.Duration, timeout time.Duration, warnMs float64) (*Ping, error) {
	p := Ping{
		BasicAggregateRoot: NewBasicAggregateRoot(),
		resourceName:       resourceName,
		target:             target,
		interval:           interval,
		timeout:            timeout,
		warnMs:             warnMs,
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Ping) ResourceName() string {
	return p.resourceName
}

func (p *Ping) Target() net.IPAddr {
	return p.target
}

func (p *Ping) Interval() time.Duration {
	return p.interval
}

func (p *Ping) Timeout() time.Duration {
	return p.timeout
}

func (p *Ping) WarnMs() float64 {
	return p.warnMs
}

func (p *Ping) validate() error {
	switch {
	case p.resourceName == "":
		return ErrInvalidPingResourceName
	case p.target.IP == nil:
		return ErrInvalidPingTarget
	case p.interval <= 0:
		return ErrInvalidPingInterval
	case p.timeout <= 0:
		return ErrInvalidPingTimeout
	case p.warnMs <= 0:
		return ErrInvalidPingWanMs
	default:
		return nil
	}
}

func (p *Ping) AddLatency(lat Latency) {
	status := ReadyStatus
	if lat.Float64() > p.warnMs {
		status = DegradedStatus
	}

	p.Record(PingEvent{
		BasicEvent: NewBasicEvent(uuid.New(), PingResultEventName, p.resourceName),
		Status:     status.String(),
		Latency:    lat.Float64(),
	})
}

func (p *Ping) AddFail() {
	p.Record(PingEvent{
		BasicEvent: NewBasicEvent(uuid.New(), PingResultEventName, p.resourceName),
		Status:     FailStatus.String(),
		Latency:    999,
	})
}
