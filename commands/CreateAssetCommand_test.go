package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/tests"
	"testing"
)

type CreateAssetCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateAssetCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateAssetCommandTestSuite))
}

func (suite *CreateAssetCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(suite.DbConn())

	assetTypeId := suite.InsertRow("asset_types", tests.AssetTypeValues(nil))
	manufacturerId := suite.InsertRow("manufacturers", tests.ManufacturerValues(nil))
	parentId := suite.InsertRow("assets", tests.AssetValues(&assetTypeId, &manufacturerId, nil, nil))

	serialNUmber := "S1234567890"
	batchNumber := "B1234567890"
	request := CreateAssetCommand{
		AssetTypeId:    assetTypeId,
		Name:           "Test Asset",
		SerialNumber:   &serialNUmber,
		BatchNumber:    &batchNumber,
		ManufacturerId: &manufacturerId,
		ParentId:       &parentId,
	}

	// act
	response, err := CreateAssetCommandHandler(request, ctx)

	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	suite.VerifyRow("assets", map[string]any{
		"id":              response.Id,
		"asset_type_id":   request.AssetTypeId,
		"name":            request.Name,
		"serial_number":   request.SerialNumber,
		"batch_number":    request.BatchNumber,
		"manufacturer_id": request.ManufacturerId,
		"parent_id":       request.ParentId,
	})
}
