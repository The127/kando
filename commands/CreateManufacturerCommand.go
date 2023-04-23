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

type CreateManufacturerResponse struct {
	Id uuid.UUID
}

func CreateManufacturerCommandHandler(command CreateManufacturerCommand, ctx context.Context) (CreateManufacturerResponse, error) {
	scope := middlewares.GetScope(ctx)

	m := ioc.Get[*mediator.Mediator](scope)
	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateManufacturerResponse{}, err
	}

	var manufacturerExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."manufacturers" where "name" = $1)`,
		command.Name).
		Scan(&manufacturerExists)
	if err != nil {
		return CreateManufacturerResponse{}, err
	}

	if manufacturerExists {
		return CreateManufacturerResponse{}, httpErrors.Conflict().WithMessage("manufacturer already exists")
	}

	var manufacturerId uuid.UUID
	err = tx.QueryRow(`insert into "public"."manufacturers" 
			("name") 
			values ($1)
			returning "id"`,
		command.Name).
		Scan(&manufacturerId)
	if err != nil {
		return CreateManufacturerResponse{}, err
	}

	manufacturerCreatedEvent := events.ManufacturerCreatedEvent{
		Id: manufacturerId,
	}
	err = mediator.SendEvent(m, manufacturerCreatedEvent, ctx)

	return CreateManufacturerResponse{
		Id: manufacturerId,
	}, nil
}
