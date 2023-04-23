package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/tests"
	"testing"
)

type CreateTagCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateTagCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTagCommandTestSuite))
}

func (suite *CreateTagCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(suite.DbConn())

	request := CreateTagCommand{
		Name: "Test Tag",
	}

	// act
	response, err := CreateTagCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	suite.VerifyRow("tags", map[string]any{
		"id":   response.Id,
		"name": request.Name,
	})
}
