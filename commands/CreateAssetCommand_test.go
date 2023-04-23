package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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

	assetTypeId := fake.AssetType(suite.DbConn(), fake.WithDefaults()).Id()
	manufacturerId := fake.Manufacturer(suite.DbConn(), fake.WithDefaults()).Id()
	parentId := fake.Asset(suite.DbConn(), &manufacturerId, &assetTypeId, nil, fake.WithDefaults()).Id()

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

	a.True(fake.AssetExists(suite.DbConn(), &manufacturerId, &assetTypeId, &parentId, fake.WithFields(
		"name", request.Name,
		"serial_number", *request.SerialNumber,
		"batch_number", *request.BatchNumber,
	).WithId(response.Id)))
}
