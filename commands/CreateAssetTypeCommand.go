package commands

import (
	"context"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateAssetTypeCommand struct {
	Name string
}

func CreateAssetTypeCommandHandler(command CreateAssetTypeCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var assetTypeExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."asset_types" where "name" = $1)`,
		command.Name).
		Scan(&assetTypeExists)
	if err != nil {
		return false, err
	}

	if assetTypeExists {
		return false, httpErrors.Conflict().WithMessage("asset type already exists")
	}

	_, err = tx.Exec(`insert into "public"."asset_types"
    			("name")
    			values ($1)`,
		command.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}
