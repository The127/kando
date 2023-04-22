package queries

import (
	"kando-backend/httpErrors"
	"net/url"
	"strconv"
)

type PaginationInfo struct {
	PageSize   int
	PageNumber int
}

type PagedResponse[TResponse any] struct {
	TotalCount int
	Data       []TResponse
}

func PagingFromQuery(query url.Values) (PaginationInfo, error) {
	pageSize, ok := query["limit"]
	if !ok {
		pageSize = []string{"10"}
	}

	parsedPageSize, err := strconv.Atoi(pageSize[0])
	if err != nil {
		return PaginationInfo{}, httpErrors.BadRequest().WithMessage("Invalid page size: " + pageSize[0])
	} else if parsedPageSize < 1 {
		parsedPageSize = 1
	} else if parsedPageSize > 100 {
		parsedPageSize = 100
	}

	pageNumber, ok := query["offset"]
	if !ok {
		pageNumber = []string{"0"}
	}

	parsedPageNumber, err := strconv.Atoi(pageNumber[0])
	if err != nil {
		return PaginationInfo{}, httpErrors.BadRequest().WithMessage("Invalid page number: " + pageNumber[0])
	} else if parsedPageNumber < 0 {
		parsedPageNumber = 0
	}

	return PaginationInfo{
		PageSize:   parsedPageSize,
		PageNumber: parsedPageNumber,
	}, nil
}
