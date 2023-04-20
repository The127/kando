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

type CreateSessionRequestDto struct {
	Username string
	Password string
}

func (dto CreateSessionRequestDto) Validate() error {
	return validator.NewFluentValidator[CreateSessionRequestDto]("dto").
		Add(validator.SubValidator[CreateSessionRequestDto, string](
			"username",
			func(dto CreateSessionRequestDto) string {
				return dto.Username
			},
			validationRules.UserUsername)).
		Validate(dto)
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateSessionRequestDto
	err := utils.DecodeJson(r.Body, &dto)
	if err != nil {
		rcs.Error(err)
		return
	}

	err = dto.Validate()
	if err != nil {
		rcs.Error(err)
		return
	}

	request := commands.CreateSessionCommand{
		Username: dto.Username,
		Password: dto.Password,
	}

	resp, err := mediator.Send[commands.CreateSessionResponse](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	http.SetCookie(w, utils.CreateSessionCookie(resp.Id))
}
