package handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

type UpdateManufacturerRequestDto struct {
	Name string
}

func (dto UpdateManufacturerRequestDto) Validate() error {
	return validator.NewFluentValidator[UpdateManufacturerRequestDto]("dto").
		Add(validator.SubValidator[UpdateManufacturerRequestDto, string](
			"name",
			func(dto UpdateManufacturerRequestDto) string {
				return dto.Name
			},
			validationRules.ManufacturerName)).
		Validate(dto)
}

func UpdateManufacturerHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	routeParams := mux.Vars(r)
	manufacturerId := uuid.MustParse(routeParams["manufacturerId"])

	var dto UpdateManufacturerRequestDto
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

	request := commands.UpdateManufacturerCommand{
		Id:   manufacturerId,
		Name: dto.Name,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
