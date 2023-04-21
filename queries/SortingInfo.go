package queries

import (
	"kando-backend/httpErrors"
	"net/url"
	"strings"
)

type SortingInfo struct {
	ColumnName string
	Direction  string
}

func (si SortingInfo) Build() string {
	return si.ColumnName + " " + si.Direction
}

func SortingFromQuery(query url.Values) ([]SortingInfo, error) {
	// sort=<name>:<dir>
	sortStrings := query["sort"]
	var result []SortingInfo

	for _, sortString := range sortStrings {
		sortingInfo := SortingInfo{}

		if sortString != "" {
			sortingInfo.ColumnName, sortingInfo.Direction = parseSortString(sortString)

			if sortingInfo.Direction != "asc" && sortingInfo.Direction != "desc" {
				return nil, httpErrors.BadRequest().
					WithMessage("Invalid sort direction: " + sortingInfo.Direction)
			}
		}

		result = append(result, sortingInfo)
	}

	return result, nil
}

func parseSortString(sortString string) (string, string) {
	sortStringParts := splitByFirst(sortString, ":")

	if len(sortStringParts) == 1 {
		return sortStringParts[0], "asc"
	}

	return sortStringParts[0], sortStringParts[1]
}

func splitByFirst(s string, sep string) []string {
	index := strings.Index(s, sep)

	if index == -1 {
		return []string{s}
	}

	return []string{s[:index], s[index+1:]}
}
