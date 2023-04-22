package queries

import (
	"net/http"
)

type QueryBase struct {
	Paging     PaginationInfo
	Sorting    []SortingInfo
	SearchText string
}

func BaseFromRequest(r *http.Request) (QueryBase, error) {
	query := r.URL.Query()

	paging, err := PagingFromQuery(query)
	if err != nil {
		return QueryBase{}, err
	}

	sorting, err := SortingFromQuery(query)
	if err != nil {
		return QueryBase{}, err
	}

	searchText := query.Get("q")

	return QueryBase{
		Paging:     paging,
		Sorting:    sorting,
		SearchText: searchText,
	}, nil
}
