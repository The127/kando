package commands

import (
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
	"kando-backend/tests"
	"testing"
)

type UpdateAssetTypeCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunUpdateAssetTypeCommandTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateAssetTypeCommandTestSuite))
}

func (s *UpdateAssetTypeCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := testContext(s.DbConn())

	assetTypeId := fake.Manufacturer(s.DbConn(), fake.WithDefaults())

	request := UpdateAssetTypeCommand{
		Id:   assetTypeId,
		Name: "Name",
	}

	// act
	_, err := UpdateAssetTypeCommandHandler(request, ctx)

	// assert
	a := s.Assert()

	a.Nil(err)
}
