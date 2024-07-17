package sockets

import (
	"time"
)

type Event struct {
	ID     uint64         `json:"id"`
	Type   EventType      `json:"type"`
	Body   map[string]any `json:"body"`
	SentAt time.Time      `json:"sent_at"`
}

func NewEvent(id uint64, t EventType, body map[string]any) Event {
	return Event{
		ID:     id,
		Type:   t,
		Body:   body,
		SentAt: time.Now(),
	}
}
