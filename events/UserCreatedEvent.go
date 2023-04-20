package events

import "github.com/google/uuid"

type UserCreatedEvent struct {
	Id uuid.UUID
}
