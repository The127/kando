package commands

import (
	"context"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateLocationCommand struct {
	Name string
}

func CreateLocationCommandHandler(command CreateLocationCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var locationExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."locations" where "name" = $1)`,
		command.Name).
		Scan(&locationExists)
	if err != nil {
		return false, err
	}

	if locationExists {
		return false, httpErrors.Conflict().WithMessage("location already exists")
	}

	_, err = tx.Exec(`insert into "public"."locations"
    			("name")
    			values ($1)`,
		command.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}
