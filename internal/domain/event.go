package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventID() uuid.UUID
	EventName() string
	CreatedAt() time.Time
	AggregateRootID() string
}

type BasicEvent struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"resourceName"`
	At                  time.Time `json:"created_at"`
	AggregateRootIDAttr string    `json:"aggregate_root_id"`
}

func (b BasicEvent) EventID() uuid.UUID {
	return b.ID
}

func (b BasicEvent) EventName() string {
	return b.Name
}

func (b BasicEvent) CreatedAt() time.Time {
	return b.At
}

func (b BasicEvent) AggregateRootID() string {
	return b.AggregateRootIDAttr
}

func NewBasicEvent(ID uuid.UUID, eventName string, aggregateRootIDAttr string) *BasicEvent {
	return &BasicEvent{ID: ID, Name: eventName, AggregateRootIDAttr: aggregateRootIDAttr, At: time.Now().UTC()}
}

type PingEvent struct {
	*BasicEvent
	Status  string  `json:"status"`
	Latency float64 `json:"latency"`
}
