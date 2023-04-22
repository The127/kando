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

type UpdateAssetTypeRequestDto struct {
	Name string
}

func (dto UpdateAssetTypeRequestDto) Validate() error {
	return validator.NewFluentValidator[UpdateAssetTypeRequestDto]("dto").
		Add(validator.SubValidator[UpdateAssetTypeRequestDto, string](
			"name",
			func(dto UpdateAssetTypeRequestDto) string {
				return dto.Name
			},
			validationRules.AssetTypeName)).
		Validate(dto)
}

func UpdateAssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto UpdateAssetTypeRequestDto
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

	request := commands.UpdateAssetTypeCommand{
		Name: dto.Name,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
