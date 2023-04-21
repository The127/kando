package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/httpErrors"
	"kando-backend/tests"
	"testing"
)

type CreateManufacturerCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateManufacturerCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateManufacturerCommandTestSuite))
}

func (s *CreateManufacturerCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := testContext(s.DbConn())

	request := CreateManufacturerCommand{
		Name: "Name",
	}

	// act
	_, err := CreateManufacturerCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)
}

func (s *CreateManufacturerCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := testContext(s.DbConn())

	request := CreateManufacturerCommand{
		Name: "Name",
	}

	// act
	_, _ = CreateManufacturerCommandHandler(request, ctx)
	_, err := CreateManufacturerCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(err, httpErrors.Conflict())
}
