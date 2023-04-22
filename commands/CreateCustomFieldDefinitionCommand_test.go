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
	ctx := testContext(suite.DbConn())

	assetTypeId := fake.AssetType(suite.DbConn(), fake.WithDefaults())

	request := CreateCustomFieldDefinitionCommand{
		AssetTypeId: assetTypeId,
		Name:        "Test Field",
		Type:        "string",
	}

	// act
	_, err := CreateCustomFieldDefinitionCommandHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)
}
