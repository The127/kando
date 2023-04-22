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
	ctx := testContext(suite.DbConn())

	request := CreateTagCommand{
		Name: "Test Tag",
	}

	// act
	_, err := CreateTagCommandHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)
}
