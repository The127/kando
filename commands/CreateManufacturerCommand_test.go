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

	s.VerifyRow("manufacturers", map[string]any{
		"id":   response.Id,
		"name": request.Name,
	})
}

func (s *CreateManufacturerCommandTestSuite) TestNameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	s.InsertRow("manufacturers", tests.ManufacturerValues(map[string]any{
		"name": "Name",
	}))

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
