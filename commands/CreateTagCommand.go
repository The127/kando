package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateTagCommand struct {
	Name string
}

type CreateTagResponse struct {
	Id uuid.UUID
}

func CreateTagCommandHandler(command CreateTagCommand, ctx context.Context) (CreateTagResponse, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateTagResponse{}, err
	}

	var tagExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."tags" where "name" = $1)`,
		command.Name).
		Scan(&tagExists)
	if err != nil {
		return CreateTagResponse{}, err
	}

	if tagExists {
		return CreateTagResponse{}, httpErrors.Conflict().WithMessage("tag already exists")
	}

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."tags"
    			("name")
    			values ($1)
    			returning id`,
		command.Name).Scan(&id)
	if err != nil {
		return CreateTagResponse{}, err
	}

	return CreateTagResponse{
		Id: id,
	}, nil
}
