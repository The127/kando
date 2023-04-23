package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func assetTypeFields(overwrites *FieldValues) *FieldValues {
	return WithFields(
		"name", faker.Word(),
	).Merge(overwrites)
}

func AssetType(db *sql.DB, overwrites *FieldValues) *FieldValues {
	fields := assetTypeFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into asset_types (name)
    									values ($1)
    									returning id`,
		get[string](fields, "name")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.WithId(id)
}

func AssetTypeExists(db *sql.DB, overwrites *FieldValues) bool {
	fields := assetTypeFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`select id from asset_types where name = $1`,
		get[string](fields, "name")).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		panic(err)
	}

	return true
}
