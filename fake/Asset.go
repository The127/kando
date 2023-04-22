package fake

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func Asset(db *sql.DB, manufacturerId *uuid.UUID, assetTypeId *uuid.UUID, parentId *uuid.UUID, fields *FieldValues) uuid.UUID {
	id := uuid.New()

	serialNumber, err := faker.RandomInt(100000, 999999)
	if err != nil {
		panic(err)
	}

	batchNumber, err := faker.RandomInt(100000, 999999)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`insert into assets (id, name, serial_number, batch_number, manufacturer_id, asset_type_id, parent_id) values ($1, $2, $3, $4, $5, $6, $7)`,
		id,
		get(fields, "name", faker.Word()),
		get(fields, "serial_number", fmt.Sprintf("S-%d", serialNumber)),
		get(fields, "batch_number", fmt.Sprintf("B-%d", batchNumber)),
		manufacturerId,
		assetTypeId,
		parentId)
	if err != nil {
		panic(err)
	}

	return id
}
