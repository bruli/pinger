package listeners

import (
	"context"
	"encoding/json"

	"github.com/bruli/pinger/internal/domain"
	"github.com/bruli/pinger/internal/infra/nats"
)

type PublishOnPingResult struct {
	publisher *nats.Publisher
}

func (p PublishOnPingResult) Listen(ctx context.Context, ev domain.Event) error {
	ping := ev.(domain.PingEvent)
	data, err := json.Marshal(&ping)
	if err != nil {
		return err
	}
	return p.publisher.Publish(ctx, data)
}

func NewPublishOnPingResult(publisher *nats.Publisher) *PublishOnPingResult {
	return &PublishOnPingResult{publisher: publisher}
}
