package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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

	a.True(fake.AssetTypeExists(s.DbConn(), fake.WithFields(
		"name", request.Name,
	).WithId(response.Id)))
}

func (s *CreateAssetTypeCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	fake.AssetType(s.DbConn(), fake.WithFields(
		"name", "Name",
	))

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
