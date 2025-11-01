package domain

import "sync"

type BasicAggregateRoot struct {
	events []Event
	sync.RWMutex
}

func NewBasicAggregateRoot() BasicAggregateRoot {
	return BasicAggregateRoot{
		events: nil,
	}
}

func (b *BasicAggregateRoot) Record(evs ...Event) {
	b.Lock()
	defer b.Unlock()
	b.events = append(b.events, evs...)
}

func (b *BasicAggregateRoot) Events() []Event {
	b.RLock()
	defer b.RUnlock()
	events := b.events
	b.clearEvents()
	return events
}

func (b *BasicAggregateRoot) clearEvents() {
	b.events = nil
}
