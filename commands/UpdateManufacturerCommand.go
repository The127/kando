package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type UpdateManufacturerCommand struct {
	Id   uuid.UUID
	Name string
}

func UpdateManufacturerHandler(command UpdateManufacturerCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)
	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(`update "public"."manufacturers" set 
                                    "name" = $1 
                                where "id" = $2`,
		command.Name, command.Id)
	if err != nil {
		return false, err
	}

	return true, nil
}
