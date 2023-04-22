package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateAssetCommand struct {
	AssetTypeId    uuid.UUID
	Name           string
	SerialNumber   *string
	BatchNumber    *string
	ManufacturerId *uuid.UUID
	ParentId       *uuid.UUID
}

type CreateAssetResponse struct {
	Id uuid.UUID
}

func CreateAssetCommandHandler(request CreateAssetCommand, ctx context.Context) (uuid.UUID, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return uuid.UUID{}, err
	}

	var assetTypeExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."asset_types" where "id" = $1)`,
		request.AssetTypeId).
		Scan(&assetTypeExists)
	if err != nil {
		return uuid.UUID{}, err
	}

	if !assetTypeExists {
		return uuid.UUID{}, httpErrors.NotFound().WithMessage("asset type does not exist")
	}

	if request.SerialNumber != nil {
		var assetExists bool
		err = tx.QueryRow(`select exists(select 1 from "public"."assets" where "serial_number" = $1 and "asset_type_id" = $2)`,
			request.SerialNumber, request.AssetTypeId).
			Scan(&assetExists)
		if err != nil {
			return uuid.UUID{}, err
		}

		if assetExists {
			return uuid.UUID{}, httpErrors.Conflict().WithMessage("asset already exists")
		}
	}

	if request.ParentId != nil {
		var assetExists bool
		err = tx.QueryRow(`select exists(select 1 from "public"."assets" where "id" = $1)`,
			request.ParentId).
			Scan(&assetExists)
		if err != nil {
			return uuid.UUID{}, err
		}

		if !assetExists {
			return uuid.UUID{}, httpErrors.NotFound().WithMessage("parent asset does not exist")
		}
	}

	var assetId uuid.UUID
	err = tx.QueryRow(`insert into "public"."assets"
    			("asset_type_id", "name", "serial_number", "batch_number", "manufacturer_id", "parent_id")
    			values ($1, $2, $3, $4, $5, $6)
    			returning "id"`,
		request.AssetTypeId, request.Name, request.SerialNumber, request.BatchNumber, request.ManufacturerId, request.ParentId).
		Scan(&assetId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return assetId, nil
}
