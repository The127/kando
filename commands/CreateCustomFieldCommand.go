package commands

import (
	"context"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateCustomFieldCommand struct {
	AssetTypeId             uuid.UUID
	CustomFieldDefinitionId uuid.UUID
	Value                   any
}

func CreateCustomFieldCommandHandler(command CreateCustomFieldCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return false, err
	}

	var customFieldExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."custom_fields" where "asset_id" = $1 and "custom_field_definition_id" = $2)`,
		command.AssetTypeId, command.CustomFieldDefinitionId).
		Scan(&customFieldExists)
	if err != nil {
		return false, err
	}

	if customFieldExists {
		return false, httpErrors.Conflict().WithMessage("custom field already exists")
	}

	customFieldValue := make(map[string]any)
	customFieldValue["value"] = command.Value

	_, err = tx.Exec(`insert into "public"."custom_fields"
				("asset_id", "custom_field_definition_id", "value")
				values ($1, $2, $3)`,
		command.AssetTypeId, command.CustomFieldDefinitionId, customFieldValue)
	if err != nil {
		return false, err
	}

	return true, nil
}
