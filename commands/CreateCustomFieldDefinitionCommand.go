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
	FieldType   string
	AssetTypeId uuid.UUID
}

type CreateCustomFieldDefinitionResponse struct {
	Id uuid.UUID
}

func CreateCustomFieldDefinitionCommandHandler(command CreateCustomFieldDefinitionCommand, ctx context.Context) (CreateCustomFieldDefinitionResponse, error) {
	scope := middlewares.GetScope(ctx)

	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateCustomFieldDefinitionResponse{}, err
	}

	var customFieldDefinitionExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."custom_field_definitions" where "name" = $1 and "asset_type_id" = $2)`,
		command.Name, command.AssetTypeId).
		Scan(&customFieldDefinitionExists)
	if err != nil {
		return CreateCustomFieldDefinitionResponse{}, err
	}

	if customFieldDefinitionExists {
		return CreateCustomFieldDefinitionResponse{}, httpErrors.Conflict().WithMessage("custom field definition already exists")
	}

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."custom_field_definitions"
    			("name", "type", "asset_type_id")
    			values ($1, $2, $3)
    			returning id`,
		command.Name, command.FieldType, command.AssetTypeId).Scan(&id)
	if err != nil {
		return CreateCustomFieldDefinitionResponse{}, err
	}

	return CreateCustomFieldDefinitionResponse{
		Id: id,
	}, nil

}
