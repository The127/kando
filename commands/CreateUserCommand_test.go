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

	s.VerifyRow("users", map[string]any{
		"id":           response.Id,
		"display_name": request.DisplayName,
		"username":     request.Username,
	})
}

func (s *CreateUserCommandTestSuite) TestUsernameAlreadyExists() {
	// arrange
	ctx := startTestContext(s.DbConn())

	s.InsertRow("users", tests.UserValues(map[string]any{
		"username": "alreadyExists",
	}))

	request := CreateUserCommand{
		DisplayName: "DisplayName",
		Username:    "alreadyExists",
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
