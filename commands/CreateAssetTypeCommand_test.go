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
	ctx := testContext(s.DbConn())

	request := CreateAssetTypeCommand{
		Name: "Name",
	}

	// act
	_, err := CreateAssetTypeCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)
}

func (s *CreateAssetTypeCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := testContext(s.DbConn())

	request := CreateAssetTypeCommand{
		Name: "Name",
	}

	// act
	_, _ = CreateAssetTypeCommandHandler(request, ctx)
	_, err := CreateAssetTypeCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Conflict(), err)
}
