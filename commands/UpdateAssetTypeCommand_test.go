package commands

import (
	"github.com/stretchr/testify/suite"
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
	ctx := startTestContext(s.DbConn())

	assetTypeId := s.InsertRow("asset_types", tests.AssetTypeValues(nil))

	request := UpdateAssetTypeCommand{
		Id:   assetTypeId,
		Name: "Name",
	}

	// act
	_, err := UpdateAssetTypeCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := s.Assert()

	a.Nil(err)

	s.VerifyRow("asset_types", map[string]any{
		"id":   request.Id,
		"name": request.Name,
	})
}
