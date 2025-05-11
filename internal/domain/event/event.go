package event

import "time"

// Event adalah interface untuk domain events
type Event interface {
	// GetType mengembalikan tipe event
	GetType() string

	// GetOccurredAt mengembalikan waktu event terjadi
	GetOccurredAt() time.Time
}

// BaseEvent adalah implementasi dasar dari Event
type BaseEvent struct {
	Type       string
	OccurredAt time.Time
}

// GetType mengembalikan tipe event
func (e BaseEvent) GetType() string {
	return e.Type
}

// GetOccurredAt mengembalikan waktu event terjadi
func (e BaseEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// NewBaseEvent membuat instance BaseEvent baru
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		Type:       eventType,
		OccurredAt: time.Now(),
	}
}
