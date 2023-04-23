package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func customFieldDefinitionFields(assetTypeId uuid.UUID, overwrites *FieldValues) *FieldValues {
	return WithFields(
		"name", faker.Word(),
		"field_type", "string",
		"asset_type_id", assetTypeId,
	).Merge(overwrites)
}

func CustomFieldDefinition(db *sql.DB, assetTypeId uuid.UUID, overwrites *FieldValues) *FieldValues {
	fields := customFieldDefinitionFields(assetTypeId, overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into custom_field_definitions (name, field_type, asset_type_id)
									values ($1, $2, $3)
									returning id`,
		get[string](fields, "name"),
		get[string](fields, "field_type"),
		get[uuid.UUID](fields, "asset_type_id")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.Merge(withId(id))
}

func CustomFieldDefinitionExists(db *sql.DB, assetTypeId uuid.UUID, overwrites *FieldValues) bool {
	fields := customFieldDefinitionFields(assetTypeId, overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from custom_field_definitions where 
									  id = $1 and
									  name = $2 and
									  field_type = $3 and
									  asset_type_id = $4)`,
		fields.Id(),
		get[string](fields, "name"),
		get[string](fields, "field_type"),
		get[uuid.UUID](fields, "asset_type_id")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
