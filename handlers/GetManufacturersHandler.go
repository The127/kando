package handlers

import (
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/middlewares"
	"kando-backend/queries"
	"kando-backend/services"
	"kando-backend/utils"
	"net/http"
)

func GetManufacturersHandler(w http.ResponseWriter, r *http.Request) {
	scope := middlewares.GetScope(r.Context())

	rcs := ioc.Get[*services.RequestContextService](scope)
	m := ioc.Get[*mediator.Mediator](scope)

	query, err := queries.BaseFromRequest(r)
	if err != nil {
		rcs.Error(err)
		return
	}

	request := queries.GetManufacturersQuery{
		QueryBase: query,
	}

	response, err := mediator.Send[queries.GetManufacturersResponse](m, request, r.Context())
	if err != nil {
		rcs.Error(err)
		return
	}

	err = utils.WriteJson(w, response)
	if err != nil {
		rcs.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
