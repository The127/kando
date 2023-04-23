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
	ctx := startTestContext(s.DbConn())

	manufacturerId := fake.Manufacturer(s.DbConn(), fake.WithDefaults()).Id()

	request := UpdateManufacturerCommand{
		Id:   manufacturerId,
		Name: "Name",
	}

	// act
	_, err := UpdateManufacturerHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := s.Assert()

	a.Nil(err)

	a.True(fake.ManufacturerExists(s.DbConn(), fake.WithFields(
		"name", request.Name,
	).WithId(request.Id)))
}
