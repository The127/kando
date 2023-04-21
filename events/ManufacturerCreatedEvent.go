package events

import "github.com/google/uuid"

type ManufacturerCreatedEvent struct {
	Id uuid.UUID
}
