package events

import "github.com/google/uuid"

type SessionCreatedEvent struct {
	Id uuid.UUID
}
