package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func Manufacturer(db *sql.DB, fields *FieldValues) uuid.UUID {
	id := uuid.New()

	_, err := db.Exec(`insert into manufacturers (id, name) values ($1, $2)`,
		id,
		get(fields, "name", faker.Word()+" Inc."))
	if err != nil {
		panic(err)
	}

	return id
}
