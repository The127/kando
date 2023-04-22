package queries

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	"kando-backend/ioc"
	"kando-backend/log"
	"kando-backend/middlewares"
)

type GetAssetTypesQuery struct {
	QueryBase
}

type GetAssetTypesResponse struct {
	Name string
}

func GetAssetTypesQueryHandler(query GetAssetTypesQuery, ctx context.Context) (PagedResponse[GetAssetTypesResponse], error) {
	scope := middlewares.GetScope(ctx)

	db := ioc.Get[*sql.DB](scope)

	sb := sqlbuilder.Select("count(*) over()", "name").
		From("asset_types")

	if query.SearchText != "" {
		sb.Where(sb.Some("name", "ilike", query.SearchText))
	}

	for _, sortingInfo := range query.Sorting {
		sb.OrderBy(sortingInfo.Build())
	}

	sb.Limit(query.Paging.PageSize).
		Offset(query.Paging.PageSize * (query.Paging.PageNumber - 1))

	sqlString, args := sb.Build()
	log.Logger.Debugf("executing sql: %s", sqlString)
	rows, err := db.Query(sqlString, args...)
	if err != nil {
		return PagedResponse[GetAssetTypesResponse]{}, err
	}
	defer rows.Close()

	var totalCount int
	var result []GetAssetTypesResponse
	for rows.Next() {
		var row GetAssetTypesResponse
		err := rows.Scan(&totalCount, &row.Name)
		if err != nil {
			return PagedResponse[GetAssetTypesResponse]{}, err
		}
		result = append(result, row)
	}

	return PagedResponse[GetAssetTypesResponse]{
		TotalCount: totalCount,
		Data:       result,
	}, nil
}
