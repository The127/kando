package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
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

	assetTypeId := fake.AssetType(suite.DbConn(), fake.WithDefaults()).Id()

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

	a.True(fake.CustomFieldDefinitionExists(suite.DbConn(), assetTypeId, fake.WithFields(
		"name", request.Name,
		"field_type", request.FieldType,
	).WithId(response.Id)))
}
