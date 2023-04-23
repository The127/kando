package commands

import (
	"github.com/stretchr/testify/suite"
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

	manufacturerId := s.InsertRow("manufacturers", tests.ManufacturerValues(nil))

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

	s.VerifyRow("manufacturers", map[string]any{
		"id":   request.Id,
		"name": request.Name,
	})
}
