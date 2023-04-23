package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func tagFields(overwrites *FieldValues) *FieldValues {
	return WithFields(
		"name", faker.Word(),
	).Merge(overwrites)
}

func Tag(db *sql.DB, overwrites *FieldValues) *FieldValues {
	fields := tagFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into tags (name)
    									values ($1)
    									returning id`,
		get[string](fields, "name")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.WithId(id)
}

func TagExists(db *sql.DB, overwrites *FieldValues) bool {
	fields := tagFields(overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from tags where 
									  id = $1 and
									  name = $2)`,
		fields.Id(),
		get[string](fields, "name")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
