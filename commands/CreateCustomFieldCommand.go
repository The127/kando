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

type CreateCustomFieldResponse struct {
	Id uuid.UUID
}

func CreateCustomFieldCommandHandler(command CreateCustomFieldCommand, ctx context.Context) (CreateCustomFieldResponse, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateCustomFieldResponse{}, err
	}

	var customFieldExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."custom_fields" where "asset_id" = $1 and "custom_field_definition_id" = $2)`,
		command.AssetTypeId, command.CustomFieldDefinitionId).
		Scan(&customFieldExists)
	if err != nil {
		return CreateCustomFieldResponse{}, err
	}

	if customFieldExists {
		return CreateCustomFieldResponse{}, httpErrors.Conflict().WithMessage("custom field already exists")
	}

	customFieldValue := make(map[string]any)
	customFieldValue["value"] = command.Value

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."custom_fields"
				("asset_id", "custom_field_definition_id", "value")
				values ($1, $2, $3)
				returning id`,
		command.AssetTypeId, command.CustomFieldDefinitionId, customFieldValue).Scan(&id)
	if err != nil {
		return CreateCustomFieldResponse{}, err
	}

	return CreateCustomFieldResponse{
		Id: id,
	}, nil
}
