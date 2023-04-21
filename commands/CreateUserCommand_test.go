package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/httpErrors"
	"kando-backend/tests"
	"testing"
)

type CreateUserCommandTestSuite struct {
	tests.DbTestSuite
}

func TestCreateUserCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserCommandTestSuite))
}

func (s *CreateUserCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := testContext(s.DbConn())

	request := CreateUserCommand{
		DisplayName: "DisplayName",
		Username:    "username",
		Password:    "abcdEFGH1234!",
	}

	// act
	_, err := CreateUserCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)
}

func (s *CreateUserCommandTestSuite) TestUsernameAlreadyExists() {
	// arrange
	ctx := testContext(s.DbConn())

	request := CreateUserCommand{
		DisplayName: "DisplayName",
		Username:    "username",
		Password:    "abcdEFGH1234!",
	}

	// act
	_, _ = CreateUserCommandHandler(request, ctx)
	_, err := CreateUserCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(err, httpErrors.Conflict())
}
