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

type CreateAssetTypeRequestDto struct {
	Name string
}

func (dto CreateAssetTypeRequestDto) Validate() error {
	return validator.NewFluentValidator[CreateAssetTypeRequestDto]("dto").
		Add(validator.SubValidator[CreateAssetTypeRequestDto, string](
			"Name",
			func(dto CreateAssetTypeRequestDto) string { return dto.Name },
			validationRules.AssetTypeName)).
		Validate(dto)
}

func CreateAssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateAssetTypeRequestDto
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

	request := commands.CreateAssetUTypeCommand{
		Name: dto.Name,
	}

	_, err = mediator.Send[any](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
