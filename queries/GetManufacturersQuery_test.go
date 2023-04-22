package queries

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
	"kando-backend/tests"
	"testing"
)

type GetManufacturersQueryTestSuite struct {
	tests.DbTestSuite
}

func TestGetManufacturersQuery(t *testing.T) {
	suite.Run(t, new(GetManufacturersQueryTestSuite))
}

func (suite *GetManufacturersQueryTestSuite) TestNoMatches() {
	// arrange
	ctx := testContext(suite.DbConn())

	request := GetManufacturersQuery{
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
	result, err := GetManufacturersQueryHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(0, result.TotalCount)
	a.Equal(0, len(result.Data))
}

func (suite *GetManufacturersQueryTestSuite) TestSingleMatch() {
	// arrange
	ctx := testContext(suite.DbConn())

	fake.Manufacturer(suite.DbConn(), fake.WithDefaults())

	request := GetManufacturersQuery{
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
	result, err := GetManufacturersQueryHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(1, result.TotalCount)
	a.Equal(1, len(result.Data))
}

func (suite *GetManufacturersQueryTestSuite) TestMultipleMatches() {
	// arrange
	ctx := testContext(suite.DbConn())

	expectedCount := 5
	for i := 0; i < expectedCount; i++ {
		fake.Manufacturer(suite.DbConn(), fake.WithDefaults())
	}

	request := GetManufacturersQuery{
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
	result, err := GetManufacturersQueryHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(expectedCount, result.TotalCount)
	a.Equal(expectedCount, len(result.Data))
}

func (suite *GetManufacturersQueryTestSuite) TestPagination() {
	// arrange
	ctx := testContext(suite.DbConn())

	expectedCount := 5
	for i := 0; i < expectedCount*2; i++ {
		fake.Manufacturer(suite.DbConn(), fake.WithDefaults())
	}

	request := GetManufacturersQuery{
		QueryBase{
			Paging: PaginationInfo{
				PageSize:   expectedCount,
				PageNumber: 0,
			},
			Sorting:    nil,
			SearchText: "",
		},
	}

	// act
	result, err := GetManufacturersQueryHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	a.Equal(expectedCount*2, result.TotalCount)
	a.Equal(expectedCount, len(result.Data))
}
