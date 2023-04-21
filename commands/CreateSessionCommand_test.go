package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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
	ctx := testContext(s.DbConn())

	fake.User(s.DbConn(), fake.WithFields(
		"username", "testuser",
		"password", "testpassword"))

	request := CreateSessionCommand{
		Username: "testuser",
		Password: "testpassword",
	}

	// act
	_, err := CreateSessionCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.Nil(err)
}

func (s *CreateSessionCommandTestSuite) TestInvalidUsername() {
	// arrange
	ctx := testContext(s.DbConn())

	fake.User(s.DbConn(), fake.WithFields(
		"username", "testuser",
		"password", "testpassword"))

	request := CreateSessionCommand{
		Username: "wronguser",
		Password: "testpassword",
	}

	// act
	_, err := CreateSessionCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Unauthorized(), err)
}

func (s *CreateSessionCommandTestSuite) TestInvalidPassword() {
	// arrange
	ctx := testContext(s.DbConn())

	fake.User(s.DbConn(), fake.WithFields(
		"username", "testuser",
		"password", "testpassword"))

	request := CreateSessionCommand{
		Username: "testuser",
		Password: "wrongPassword",
	}

	// act
	_, err := CreateSessionCommandHandler(request, ctx)

	// assert
	a := assert.New(s.T())

	a.NotNil(err)
	a.IsType(httpErrors.Unauthorized(), err)
}
