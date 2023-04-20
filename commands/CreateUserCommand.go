package commands

import (
	"context"
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

func CreateUserCommandHandler(command CreateUserCommand, ctx context.Context) (any, error) {
	scope := middlewares.GetScope(ctx)

	m := ioc.Get[*mediator.Mediator](scope)
	rcs := ioc.Get[*services.RequestContextService](scope)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.MinCost)
	if err != nil {
		return false, err
	}

	tx, err := rcs.BeginTx()
	if err != nil {
		return false, err
	}

	var userExists bool
	err = tx.QueryRow(`select exists(select 1 from "public"."users" where "username" = $1)`,
		command.Username).
		Scan(&userExists)
	if err != nil {
		return false, err
	}

	if userExists {
		return false, httpErrors.Conflict().WithMessage("username is already taken")
	}

	_, err = tx.Exec(`insert into "public"."users" 
			("display_name", "username", "hashed_password")
			values ($1, $2, $3);`,
		command.DisplayName, command.Username, hashedPassword)
	if err != nil {
		return false, err
	}

	userCreatedEvent := events.UserCreatedEvent{}
	err = mediator.SendEvent(m, userCreatedEvent, ctx)
	if err != nil {
		return false, err
	}

	err = rcs.CommitTx()
	if err != nil {
		return false, err
	}

	return true, nil
}
