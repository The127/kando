package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func locationFields(overwrites *FieldValues) *FieldValues {
	return WithFields(
		"name", faker.Word()+" Town",
	).Merge(overwrites)
}

func Location(db *sql.DB, overwrites *FieldValues) *FieldValues {
	fields := locationFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into locations (name)
    									values ($1)
    									returning id`,
		get[string](fields, "name")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.WithId(id)
}

func LocationExists(db *sql.DB, overwrites *FieldValues) bool {
	fields := locationFields(overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from locations where 
									  id = $1 and
									  name = $2)`,
		fields.Id(),
		get[string](fields, "name")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
