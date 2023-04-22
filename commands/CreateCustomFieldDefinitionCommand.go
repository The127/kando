package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateCustomFieldDefinitionCommand struct {
	Name        string
	Type        string
	AssetTypeId uuid.UUID
}

func CreateCustomFieldDefinitionCommandHandler(command CreateCustomFieldDefinitionCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var customFieldDefinitionExists bool
	err = tx.QueryRow(`select exists(select exists(select 1 from "public"."custom_field_definitions" where "name" = $1 and "asset_type_id" = $2))`,
		command.Name, command.AssetTypeId).
		Scan(&customFieldDefinitionExists)
	if err != nil {
		return false, err
	}

	if customFieldDefinitionExists {
		return false, httpErrors.Conflict().WithMessage("custom field definition already exists")
	}

	_, err = tx.Exec(`insert into "public"."custom_field_definitions"
    			("name", "type", "asset_type_id")
    			values ($1, $2, $3)`,
		command.Name, command.Type, command.AssetTypeId)
	if err != nil {
		return false, err
	}

	return true, nil

}
