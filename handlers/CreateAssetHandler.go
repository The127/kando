package handlers

import (
	"github.com/google/uuid"
	"kando-backend/commands"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/services"
	"kando-backend/utils"
	"net/http"
)

type CreateAssetRequestDto struct {
	AssetTypeId    uuid.UUID
	Name           string
	SerialNumber   *string
	BatchNumber    *string
	ManufacturerId *uuid.UUID
	ParentId       *uuid.UUID
}

type CreateAssetResponseDto struct {
	Id uuid.UUID
}

func CreateAssetHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	var dto CreateAssetRequestDto
	err := utils.ReadJson(r.Body, &dto)
	if err != nil {
		rcs.Error(err)
		return
	}

	request := commands.CreateAssetCommand{
		AssetTypeId:    dto.AssetTypeId,
		Name:           dto.Name,
		SerialNumber:   dto.SerialNumber,
		BatchNumber:    dto.BatchNumber,
		ManufacturerId: dto.ManufacturerId,
		ParentId:       dto.ParentId,
	}

	response, err := mediator.Send[commands.CreateAssetResponse](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	responseDto := CreateAssetResponseDto{
		Id: response.Id,
	}
	err = utils.WriteJson(w, responseDto)
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
