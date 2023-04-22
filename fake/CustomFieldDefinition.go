package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func CustomFieldDefinition(db *sql.DB, assetTypeId uuid.UUID, fields *FieldValues) uuid.UUID {
	id := uuid.New()

	_, err := db.Exec(`insert into custom_field_definitions (id, asset_type_id, name, field_type) VALUES ($1, $2, $3, $4)`,
		id,
		assetTypeId,
		get(fields, "name", faker.Word()),
		get(fields, "field_type", "string"))
	if err != nil {
		panic(err)
	}

	return id
}
