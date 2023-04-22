package commands

import (
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
	"kando-backend/tests"
)

type UpdateManufacturerCommandTestSuite struct {
	tests.DbTestSuite
}

func (s *UpdateManufacturerCommandTestSuite) TestRunUpdateManufacturerCommandTestSuite() {
	suite.Run(s.T(), new(UpdateManufacturerCommandTestSuite))
}

func (s *UpdateManufacturerCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := testContext(s.DbConn())

	manufacturerId := fake.Manufacturer(s.DbConn(), fake.WithDefaults())

	request := UpdateManufacturerCommand{
		Id:   manufacturerId,
		Name: "Name",
	}

	// act
	_, err := UpdateManufacturerHandler(request, ctx)

	// assert
	a := s.Assert()

	a.Nil(err)
}
