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

type CreateManufacturerRequestDto struct {
	Name string
}

func (dto CreateManufacturerRequestDto) Validate() error {
	return validator.NewFluentValidator[CreateManufacturerRequestDto]("dto").
		Add(validator.SubValidator[CreateManufacturerRequestDto, string](
			"name",
			func(dto CreateManufacturerRequestDto) string {
				return dto.Name
			},
			validationRules.ManufacturerName)).
		Validate(dto)
}

func CreateManufacturerHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateManufacturerRequestDto
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

	request := commands.CreateManufacturerCommand{
		Name: dto.Name,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
