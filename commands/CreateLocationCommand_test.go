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
	ctx := testContext(suite.DbConn())

	request := CreateLocationCommand{
		Name: "Test Location",
	}

	// act
	_, err := CreateLocationCommandHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)
}
