package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/httpErrors"
	"kando-backend/tests"
	"testing"
)

type CreateAssetTypeCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateAssetTypeCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateAssetTypeCommandTestSuite))
}

func (s *CreateAssetTypeCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(s.DbConn())

	request := CreateAssetTypeCommand{
		Name: "Name",
	}

	// act
	response, err := CreateAssetTypeCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)

	s.VerifyRow("asset_types", map[string]interface{}{
		"id":   response.Id,
		"name": request.Name,
	})
}

func (s *CreateAssetTypeCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	s.InsertRow("asset_types", tests.AssetTypeValues(map[string]any{
		"name": "Name",
	}))

	request := CreateAssetTypeCommand{
		Name: "Name",
	}

	// act
	_, err := CreateAssetTypeCommandHandler(request, ctx)
	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Conflict(), err)
}
