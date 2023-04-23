package commands

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"kando-backend/events"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/services"
)

type CreateUserCommand struct {
	DisplayName string
	Username    string
	Password    string
}

type CreateUserResponse struct {
	Id uuid.UUID
}

func CreateUserCommandHandler(command CreateUserCommand, ctx context.Context) (CreateUserResponse, error) {
	scope := middlewares.GetScope(ctx)

	m := ioc.Get[*mediator.Mediator](scope)
	rcs := ioc.Get[*services.RequestContextService](scope)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.MinCost)
	if err != nil {
		return CreateUserResponse{}, err
	}

	tx, err := rcs.GetTx()
	if err != nil {
		return CreateUserResponse{}, err
	}

	var userExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."users" where "username" = $1)`,
		command.Username).
		Scan(&userExists)
	if err != nil {
		return CreateUserResponse{}, err
	}

	if userExists {
		return CreateUserResponse{}, httpErrors.Conflict().WithMessage("username is already taken")
	}

	var id uuid.UUID
	err = tx.QueryRow(`insert into "public"."users" 
			("display_name", "username", "hashed_password")
			values ($1, $2, $3)
			returning "id";`,
		command.DisplayName, command.Username, hashedPassword).Scan(&id)
	if err != nil {
		return CreateUserResponse{}, err
	}

	userCreatedEvent := events.UserCreatedEvent{
		Id: id,
	}
	err = mediator.SendEvent(m, userCreatedEvent, ctx)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return CreateUserResponse{
		Id: id,
	}, nil
}
