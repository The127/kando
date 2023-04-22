package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func AssetType(db *sql.DB, fields *FieldValues) uuid.UUID {
	id := uuid.New()

	_, err := db.Exec(`insert into asset_types (id, name, manufacturer_id) VALUES ($1, $2, $3)`,
		id,
		get(fields, "name", faker.Word()))
	if err != nil {
		panic(err)
	}

	return id
}
