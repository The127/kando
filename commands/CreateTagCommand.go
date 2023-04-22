package commands

import (
	"context"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateTagCommand struct {
	Name string
}

func CreateTagCommandHandler(command CreateTagCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var tagExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."tags" where "name" = $1)`,
		command.Name).
		Scan(&tagExists)
	if err != nil {
		return false, err
	}

	if tagExists {
		return false, httpErrors.Conflict().WithMessage("tag already exists")
	}

	_, err = tx.Exec(`insert into "public"."tags"
    			("name")
    			values ($1)`,
		command.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}
