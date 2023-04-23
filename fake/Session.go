package fake

import (
	"database/sql"
	"github.com/google/uuid"
)

func sessionFields(userId uuid.UUID, overwrites *FieldValues) *FieldValues {
	return WithFields(
		"user_id", userId,
	).Merge(overwrites)
}

func Session(db *sql.DB, userId uuid.UUID, overwrites *FieldValues) *FieldValues {
	fields := sessionFields(userId, overwrites)

	var id uuid.UUID
	err := db.QueryRow(`insert into sessions (user_id)
    										values ($1)
    										returning id`,
		get[uuid.UUID](fields, "user_id")).Scan(&id)
	if err != nil {
		panic(err)
	}

	return fields.WithId(id)
}

func SessionExists(db *sql.DB, userId uuid.UUID, overwrites *FieldValues) bool {
	fields := sessionFields(userId, overwrites)

	var exists bool
	err := db.QueryRow(`select exists(select 1 from sessions where 
									  id = $1 and
									  user_id = $2)`,
		fields.Id(),
		get[uuid.UUID](fields, "user_id")).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}
