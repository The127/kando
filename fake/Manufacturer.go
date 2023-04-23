package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func manufacturerFields(overwrites *FieldValues) *FieldValues {
	return WithFields(
		"name", faker.Word()+" Inc.",
	).Merge(overwrites)
}

func Manufacturer(db *sql.DB, overwrites *FieldValues) *FieldValues {
	fields := manufacturerFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into manufacturers (name)
    								values ($1)
    								returning id`,
		get[string](fields, "name")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.Merge(withId(id))
}

func ManufacturerExists(db *sql.DB, overwrites *FieldValues) bool {
	fields := manufacturerFields(overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from manufacturers where 
									  id = $1 and
									  name = $2)`,
		fields.Id(),
		get[string](fields, "name")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
