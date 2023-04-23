package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateLocationCommand struct {
	Name string
}

type CreateLocationResponse struct {
	Id uuid.UUID
}

func CreateLocationCommandHandler(command CreateLocationCommand, ctx context.Context) (CreateLocationResponse, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateLocationResponse{}, err
	}

	var locationExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."locations" where "name" = $1)`,
		command.Name).
		Scan(&locationExists)
	if err != nil {
		return CreateLocationResponse{}, err
	}

	if locationExists {
		return CreateLocationResponse{}, httpErrors.Conflict().WithMessage("location already exists")
	}

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."locations"
    			("name")
    			values ($1)
    			returning id`,
		command.Name).Scan(&id)
	if err != nil {
		return CreateLocationResponse{}, err
	}

	return CreateLocationResponse{
		Id: id,
	}, nil
}
