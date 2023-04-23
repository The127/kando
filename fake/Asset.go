package fake

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func assetFields(overwrites *FieldValues) *FieldValues {
	serialNumber, err := faker.RandomInt(100000, 999999)
	if err != nil {
		panic(err)
	}

	batchNumber, err := faker.RandomInt(100000, 999999)
	if err != nil {
		panic(err)
	}

	return WithFields(
		"name", faker.Word(),
		"serial_number", fmt.Sprintf("S-%d", serialNumber),
		"batch_number", fmt.Sprintf("B-%d", batchNumber),
	).Merge(overwrites)
}

func Asset(db *sql.DB, manufacturerId *uuid.UUID, assetTypeId *uuid.UUID, parentId *uuid.UUID, overwrites *FieldValues) *FieldValues {
	fields := assetFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into assets (name, serial_number, batch_number, manufacturer_id, asset_type_id, parent_id)
    									values ($1, $2, $3, $4, $5, $6)
    									returning id`,
		get[string](fields, "name"),
		get[string](fields, "serial_number"),
		get[string](fields, "batch_number"),
		manufacturerId,
		assetTypeId,
		parentId).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.WithId(id)
}

func AssetExists(db *sql.DB, manufacturerId *uuid.UUID, assetTypeId *uuid.UUID, parentId *uuid.UUID, overwrites *FieldValues) bool {
	fields := assetFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`select id from assets where name = $1 and serial_number = $2 and batch_number = $3 and manufacturer_id = $4 and asset_type_id = $5 and parent_id = $6`,
		get[string](fields, "name"),
		get[string](fields, "serial_number"),
		get[string](fields, "batch_number"),
		manufacturerId,
		assetTypeId,
		parentId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}

	return true
}
