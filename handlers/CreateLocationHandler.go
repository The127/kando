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

type CreateLocationRequestDto struct {
	Name string
}

func (dto CreateLocationRequestDto) Validate() error {
	return validator.NewFluentValidator[CreateLocationRequestDto]("dto").
		Add(validator.SubValidator[CreateLocationRequestDto, string](
			"Name",
			func(dto CreateLocationRequestDto) string { return dto.Name },
			validationRules.LocationName)).
		Validate(dto)
}

func CreateLocationHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateLocationRequestDto
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

	request := commands.CreateLocationCommand{
		Name: dto.Name,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
