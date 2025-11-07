package nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

const PingSubjet = "ping.created"

type Publisher struct {
	conn *nats.Conn
}

func (p Publisher) Close() {
	_ = p.conn.Close
}

func (p Publisher) Publish(ctx context.Context, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return p.conn.Publish(PingSubjet, data)
	}
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Publisher{conn: conn}, nil
}
