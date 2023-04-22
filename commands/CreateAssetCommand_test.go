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
	ctx := testContext(suite.DbConn())

	assetTypeId := fake.AssetType(suite.DbConn(), fake.WithDefaults())
	manufacturerId := fake.Manufacturer(suite.DbConn(), fake.WithDefaults())
	parentId := fake.Asset(suite.DbConn(), &manufacturerId, &assetTypeId, nil, fake.WithDefaults())

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
	_, err := CreateAssetCommandHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)
}
