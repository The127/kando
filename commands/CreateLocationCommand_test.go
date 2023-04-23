package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/tests"
	"testing"
)

type CreateLocationCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateLocationCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateLocationCommandTestSuite))
}

func (suite *CreateLocationCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(suite.DbConn())

	request := CreateLocationCommand{
		Name: "Test Location",
	}

	// act
	response, err := CreateLocationCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	suite.VerifyRow("locations", map[string]any{
		"id":   response.Id,
		"name": request.Name,
	})
}
