package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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
	ctx := startTestContext(s.DbConn())

	request := CreateManufacturerCommand{
		Name: "Name",
	}

	// act
	response, err := CreateManufacturerCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)

	a.True(fake.ManufacturerExists(s.DbConn(), fake.WithFields(
		"name", request.Name,
	).WithId(response.Id)))
}

func (s *CreateManufacturerCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	request := CreateManufacturerCommand{
		Name: "Name",
	}

	// act
	_, _ = CreateManufacturerCommandHandler(request, ctx)
	_, err := CreateManufacturerCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Conflict(), err)
}
