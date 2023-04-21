package commands

import (
	"github.com/google/uuid"
	"kando-backend/ioc"
	"kando-backend/services"
)

type DeleteSessionCommand struct {
	SessionId uuid.UUID
}

func DeleteSessionCommandHandler(command DeleteSessionCommand, scope *ioc.DependencyProvider) (any, error) {
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

	if err != nil {
		return false, err
	}

	return true, nil
}
