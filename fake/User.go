package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func User(db *sql.DB, fields *FieldValues) uuid.UUID {
	id := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(get(fields, "password", "password")), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`insert into users (id, display_name, username, hashed_password) values ($1, $2, $3, $4)`,
		id,
		get(fields, "display_name", faker.Name()),
		get(fields, "username", faker.Username()),
		hashedPassword)
	if err != nil {
		panic(err)
	}

	return id
}
