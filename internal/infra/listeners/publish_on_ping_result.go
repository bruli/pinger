package listeners

import (
	"context"

	"github.com/bruli/pinger/internal/domain"
	"github.com/bruli/pinger/internal/infra/nats"
	"github.com/bruli/pinger/pkg/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PublishOnPingResult struct {
	publisher *nats.Publisher
}

func (p PublishOnPingResult) Listen(ctx context.Context, ev domain.Event) error {
	ping := ev.(domain.PingEvent)
	event := events.PingResult{
		Resource:  ping.AggregateRootID(),
		Status:    ping.Status,
		Latency:   float32(ping.Latency),
		CreatedAt: timestamppb.New(ping.CreatedAt()),
	}
	data, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	return p.publisher.Publish(ctx, data)
}

func NewPublishOnPingResult(publisher *nats.Publisher) *PublishOnPingResult {
	return &PublishOnPingResult{publisher: publisher}
}
