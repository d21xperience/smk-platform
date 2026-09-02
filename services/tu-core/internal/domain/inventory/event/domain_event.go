package event

import "time"

type DomainEvent interface {
	GetEventName() string
	GetOccurredAt() time.Time
	GetAggregateID() string
	GetPayload() interface{}
}
