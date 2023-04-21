package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/events"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateManufacturerCommand struct {
	Name string
}

func CreateManufacturerCommandHandler(command CreateManufacturerCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	m := ioc.Get[*mediator.Mediator](scope)
	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var manufacturerExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."manufacturers" where "name" = $1)`,
		command.Name).
		Scan(&manufacturerExists)
	if err != nil {
		return false, err
	}

	if manufacturerExists {
		return false, httpErrors.Conflict().WithMessage("manufacturer already exists")
	}

	var manufacturerId uuid.UUID
	err = tx.QueryRow(`insert into "public"."manufacturers" 
			("name") 
			values ($1)
			returning "id"`,
		command.Name).
		Scan(&manufacturerId)
	if err != nil {
		return false, err
	}

	manufacturerCreatedEvent := events.ManufacturerCreatedEvent{
		Id: manufacturerId,
	}
	err = mediator.SendEvent(m, manufacturerCreatedEvent, ctx)

	return true, nil
}
