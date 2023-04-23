package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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
	ctx := startTestContext(s.DbConn())

	request := CreateUserCommand{
		DisplayName: "DisplayName",
		Username:    "username",
		Password:    "abcdEFGH1234!",
	}

	// act
	response, err := CreateUserCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)

	a.True(fake.UserExists(s.DbConn(), fake.WithFields(
		"username", request.Username,
		"display_name", request.DisplayName,
	).WithId(response.Id)))
}

func (s *CreateUserCommandTestSuite) TestUsernameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	fake.User(s.DbConn(), fake.WithFields(
		"username", "username"))

	request := CreateUserCommand{
		DisplayName: "DisplayName",
		Username:    "username",
		Password:    "abcdEFGH1234!",
	}

	// act
	_, err := CreateUserCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Conflict(), err)
}
