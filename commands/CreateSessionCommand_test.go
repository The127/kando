package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/httpErrors"
	"kando-backend/tests"
	"testing"
)

type CreateSessionCommandTestSuite struct {
	tests.DbTestSuite
}

func TestCreateSessionCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateSessionCommandTestSuite))
}

func (s *CreateSessionCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(s.DbConn())

	userId := s.InsertRow("users", tests.UserValues(map[string]any{
		"username": "testuser",
		"password": "testpassword",
	}))

	request := CreateSessionCommand{
		Username: "testuser",
		Password: "testpassword",
	}

	// act
	response, err := CreateSessionCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)

	s.VerifyRow("sessions", map[string]any{
		"id":      response.Id,
		"user_id": userId,
	})
}

func (s *CreateSessionCommandTestSuite) TestInvalidUsername() {
	// arrange
	ctx := startTestContext(s.DbConn())

	s.InsertRow("users", tests.UserValues(map[string]any{
		"username": "testuser",
	}))

	request := CreateSessionCommand{
		Username: "wronguser",
		Password: "testpassword",
	}

	// act
	_, err := CreateSessionCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Unauthorized(), err)
}

func (s *CreateSessionCommandTestSuite) TestInvalidPassword() {
	// arrange
	ctx := startTestContext(s.DbConn())

	s.InsertRow("users", tests.UserValues(map[string]any{
		"username": "testuser",
		"password": "testpassword",
	}))

	request := CreateSessionCommand{
		Username: "testuser",
		Password: "wrongPassword",
	}

	// act
	_, err := CreateSessionCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Unauthorized(), err)
}
