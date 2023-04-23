package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/tests"
	"testing"
)

type CreateCustomFieldDefinitionCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateCustomFieldDefinitionCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCustomFieldDefinitionCommandTestSuite))
}

func (suite *CreateCustomFieldDefinitionCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := startTestContext(suite.DbConn())

	assetTypeId := suite.InsertRow("asset_types", tests.AssetTypeValues(nil))

	request := CreateCustomFieldDefinitionCommand{
		AssetTypeId: assetTypeId,
		Name:        "Test Field",
		FieldType:   "text",
	}

	// act
	response, err := CreateCustomFieldDefinitionCommandHandler(request, ctx)
	closeTestContext(ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)

	suite.VerifyRow("custom_field_definitions", map[string]any{
		"id":            response.Id,
		"asset_type_id": request.AssetTypeId,
		"name":          request.Name,
		"field_type":    request.FieldType,
	})
}
