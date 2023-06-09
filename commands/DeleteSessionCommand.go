package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type DeleteSessionCommand struct {
	SessionId uuid.UUID
}

func DeleteSessionCommandHandler(command DeleteSessionCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)
	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(`delete from "public"."sesions" where "id" = $1`,
		command.SessionId)
	if err != nil {
		return false, err
	}

	return true, nil
}
