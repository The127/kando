package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateAssetTypeCommand struct {
	Name string
}

type CreateAssetTypeResponse struct {
	Id uuid.UUID
}

func CreateAssetTypeCommandHandler(command CreateAssetTypeCommand, ctx context.Context) (CreateAssetTypeResponse, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateAssetTypeResponse{}, err
	}

	var assetTypeExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."asset_types" where "name" = $1)`,
		command.Name).
		Scan(&assetTypeExists)
	if err != nil {
		return CreateAssetTypeResponse{}, err
	}

	if assetTypeExists {
		return CreateAssetTypeResponse{}, httpErrors.Conflict().WithMessage("asset type already exists")
	}

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."asset_types"
    			("name")
    			values ($1)
    			returning id`,
		command.Name).Scan(&id)
	if err != nil {
		return CreateAssetTypeResponse{}, err
	}

	return CreateAssetTypeResponse{
		Id: id,
	}, nil
}
