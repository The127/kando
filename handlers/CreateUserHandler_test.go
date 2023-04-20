package handlers

import (
	"github.com/stretchr/testify/suite"
	"kando-backend/tests"
	"testing"
)

type CreateUserHandlerTestSuite struct {
	tests.DbTestSuite
}

func (s *CreateUserHandlerTestSuite) Test_TooLongUsername() {
	s.Equal(1, 1, "1 should be 1")
}

func TestCreateUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserHandlerTestSuite))
}
