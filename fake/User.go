package fake

import (
	"database/sql"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func userFields(overwrites *FieldValues) *FieldValues {
	fields := WithFields(
		"display_name", faker.Name(),
		"username", faker.Username(),
		"password", faker.Password(),
	).Merge(overwrites)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(get[string](overwrites, "password")), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	fields = fields.Merge(WithFields(
		"hashed_password", hashedPassword,
	))

	return fields
}

func User(db *sql.DB, overwrites *FieldValues) *FieldValues {
	fields := userFields(overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into users (display_name, username, hashed_password) 
								values ($1, $2, $3) 
								returning id`,
		get[string](fields, "display_name"),
		get[string](fields, "username"),
		get[[]byte](fields, "hashed_password")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.Merge(withId(id))
}

func UserExists(db *sql.DB, overwrites *FieldValues) bool {
	fields := userFields(overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from users where 
                                      id = $1 and
                                      display_name = $2 and
                                      username = $3 and
                                      hashed_password = $4)`,
		fields.Id(),
		get[string](fields, "display_name"),
		get[string](fields, "username"),
		get[[]byte](fields, "hashed_password")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
