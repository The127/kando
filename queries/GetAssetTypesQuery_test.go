package queries

import (
	"github.com/stretchr/testify/assert"
	"kando-backend/fake"
	"kando-backend/tests"
)

type GetAssetTypesQueryTestSuite struct {
	tests.DbTestSuite
}

func (suite *GetAssetTypesQueryTestSuite) TestNoMatches() {
	// arrange
	ctx := testContext(suite.DbConn())

	request := GetAssetTypesQuery{
		QueryBase{
			Paging: PaginationInfo{
				PageSize:   10,
				PageNumber: 0,
			},
			Sorting:    nil,
			SearchText: "",
		},
	}

	// act
	result, err := GetAssetTypesQueryHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(0, result.TotalCount)
	a.Equal(0, len(result.Data))
}

func (suite *GetAssetTypesQueryTestSuite) TestSingleMatch() {
	// arrange
	ctx := testContext(suite.DbConn())

	fake.AssetType(suite.DbConn(), fake.WithDefaults())

	request := GetAssetTypesQuery{
		QueryBase{
			Paging: PaginationInfo{
				PageSize:   10,
				PageNumber: 0,
			},
			Sorting:    nil,
			SearchText: "",
		},
	}

	// act
	result, err := GetAssetTypesQueryHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(1, result.TotalCount)
	a.Equal(1, len(result.Data))
}

func (suite *GetAssetTypesQueryTestSuite) TestMultipleMatches() {
	// arrange
	ctx := testContext(suite.DbConn())

	expectedCount := 5
	for i := 0; i < expectedCount*2; i++ {
		fake.AssetType(suite.DbConn(), fake.WithDefaults())
	}

	request := GetAssetTypesQuery{
		QueryBase{
			Paging: PaginationInfo{
				PageSize:   10,
				PageNumber: 0,
			},
			Sorting:    nil,
			SearchText: "",
		},
	}

	// act
	result, err := GetAssetTypesQueryHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(3, result.TotalCount)
	a.Equal(3, len(result.Data))
}
