package commands

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"kando-backend/events"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateSessionCommand struct {
	Username string
	Password string
}

type CreateSessionResponse struct {
	Id uuid.UUID
}

func CreateSessionCommandHandler(command CreateSessionCommand, ctx context.Context) (CreateSessionResponse, error) {
	scope := middlewares.GetScope(ctx)

	m := ioc.Get[*mediator.Mediator](scope)
	rcs := ioc.Get[*services.RequestContextService](scope)

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateSessionResponse{}, err
	}

	var userId uuid.UUID
	var hashedPassword []byte
	err = tx.QueryRow(`select "id", "hashed_password" from "public"."users" where username = $1`,
		command.Username).
		Scan(&userId, &hashedPassword)
	if err == sql.ErrNoRows {
		return CreateSessionResponse{}, httpErrors.Unauthorized()
	}
	if err != nil {
		return CreateSessionResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(command.Password))
	if err != nil {
		return CreateSessionResponse{}, httpErrors.Unauthorized()
	}

	var sessionId uuid.UUID
	err = tx.QueryRow(`insert into "public"."sessions"
			("user_id")
			values ($1)
			returning "id";`,
		userId).
		Scan(&sessionId)
	if err != nil {
		return CreateSessionResponse{}, err
	}

	sessionCreatedEvent := events.SessionCreatedEvent{
		Id: sessionId,
	}
	err = mediator.SendEvent(m, sessionCreatedEvent, ctx)
	if err != nil {
		return CreateSessionResponse{}, err
	}

	return CreateSessionResponse{
		Id: sessionId,
	}, nil
}
