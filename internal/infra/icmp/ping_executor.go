package icmp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bruli/pinger/internal/domain"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PingExecutor struct{}

func NewPingExecutor() *PingExecutor {
	return &PingExecutor{}
}

func (p2 PingExecutor) Execute(ctx context.Context, p *domain.Ping) (domain.Latency, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("error listening icmp: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("ping-test"),
		},
	}
	data, err := msg.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("error marshaling icmp message: %w", err)
	}
	start := time.Now()
	target := p.Target()
	if _, err = conn.WriteTo(data, &target); err != nil {
		return 0, fmt.Errorf("error writing icmp message: %w", err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(p.Timeout()))
	reply := make([]byte, 1500)
	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return 0, fmt.Errorf("error reading icmp message: %w", err)
	}

	elapsed := time.Since(start)
	rm, err := icmp.ParseMessage(1, reply[:n]) // 1 = ICMPv4
	if err != nil {
		return 0, fmt.Errorf("error parsing icmp message: %w", err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return domain.Latency(elapsed.Milliseconds()), nil
	case ipv4.ICMPTypeDestinationUnreachable:
		return 0, errors.New("destination unreachable")
	}
	return 0, nil
}
