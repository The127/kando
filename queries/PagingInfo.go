package queries

type PagingInfo struct {
	PageSize   int
	PageNumber int
}

type PagedResponse[TResponse any] struct {
	TotalCount int
	Data       []TResponse
}
