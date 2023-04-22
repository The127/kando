package commands

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"kando-backend/fake"
	"kando-backend/tests"
	"testing"
)

type CreateCustomFieldCommandTestSuite struct {
	tests.DbTestSuite
}

func TestRunCreateCustomFieldCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCustomFieldCommandTestSuite))
}

func (suite *CreateCustomFieldCommandTestSuite) TestValidInputs() {
	// arrange
	ctx := testContext(suite.DbConn())

	assetTypeId := fake.AssetType(suite.DbConn(), fake.WithDefaults())
	fake.CustomFieldDefinition(suite.DbConn(), assetTypeId, fake.WithDefaults())

	request := CreateCustomFieldCommand{
		AssetTypeId:             uuid.UUID{},
		CustomFieldDefinitionId: uuid.UUID{},
		Value:                   nil,
	}

	// act
	_, err := CreateCustomFieldCommandHandler(request, ctx)

	// assert
	a := assert.New(suite.T())

	a.Nil(err)
}
