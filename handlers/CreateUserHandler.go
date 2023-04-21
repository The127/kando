package handlers

import (
	"kando-backend/commands"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/services"
	"kando-backend/utils"
	"kando-backend/validationRules"
	"kando-backend/validator"
	"net/http"
)

type CreateUserRequestDto struct {
	DisplayName string
	Username    string
	Password    string
}

func (dto CreateUserRequestDto) Validate() error {
	return validator.NewFluentValidator[CreateUserRequestDto]("dto").
		Add(validator.SubValidator[CreateUserRequestDto, string](
			"displayName",
			func(dto CreateUserRequestDto) string {
				return dto.DisplayName
			},
			validationRules.UserDisplayName)).
		Add(validator.SubValidator[CreateUserRequestDto, string](
			"username",
			func(dto CreateUserRequestDto) string {
				return dto.Username
			},
			validationRules.UserUsername)).
		Validate(dto)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateUserRequestDto
	err := utils.ReadJson(r.Body, &dto)
	if err != nil {
		rcs.Error(err)
		return
	}

	err = dto.Validate()
	if err != nil {
		rcs.Error(err)
		return
	}

	request := commands.CreateUserCommand{
		DisplayName: dto.DisplayName,
		Username:    dto.Username,
		Password:    dto.Password,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
