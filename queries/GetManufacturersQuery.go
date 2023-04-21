package queries

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	"kando-backend/ioc"
	"kando-backend/middlewares"
)

type GetManufacturersQuery struct {
	Paging     PagingInfo
	Sorting    []SortingInfo
	SearchText string
}

type GetManufacturersResponse struct {
	Name string
}

func GetManufacturersQueryHandler(query GetManufacturersQuery, ctx context.Context) (PagedResponse[GetManufacturersResponse], error) {
	scope := middlewares.GetScope(ctx)

	db := ioc.Get[*sql.DB](scope)

	sb := sqlbuilder.Select("count(*) over()", "name").
		From("manufacturers")

	if query.SearchText != "" {
		sb.Where(sb.Some("name", "ilike", query.SearchText))
	}

	for _, sortingInfo := range query.Sorting {
		sb.OrderBy(sortingInfo.Build())
	}

	sb.Limit(query.Paging.PageSize).
		Offset(query.Paging.PageSize * (query.Paging.PageNumber - 1))

	rows, err := db.Query(sb.Build())
	if err != nil {
		return PagedResponse[GetManufacturersResponse]{}, err
	}
	defer rows.Close()

	var totalCount int
	var result []GetManufacturersResponse
	for rows.Next() {
		var row GetManufacturersResponse
		err := rows.Scan(&totalCount, &row.Name)
		if err != nil {
			return PagedResponse[GetManufacturersResponse]{}, err
		}
		result = append(result, row)
	}

	return PagedResponse[GetManufacturersResponse]{
		TotalCount: totalCount,
		Data:       result,
	}, nil
}
