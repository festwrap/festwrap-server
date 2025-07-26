package event

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	PlaylistCreated EventType = "playlist_created"
)

type Event interface {
	Type() EventType
}

type EventWrapper[T Event] struct {
	EventID   string    `json:"id"`
	Timestamp int64     `json:"timestamp"`
	EventType EventType `json:"type"`
	Event     T         `json:"payload"`
}

func NewEventWrapper[T Event](event T) EventWrapper[T] {
	return EventWrapper[T]{
		Timestamp: time.Now().UnixMilli(),
		EventID:   uuid.NewString(),
		Event:     event,
		EventType: event.Type(),
	}
}
