package model

import (
	"testing"

	"github.com/auto-staging/builder/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseModel(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)

	model := NewDatabaseModel(svc)

	assert.NotEmpty(t, model, "Expected model not to be empty")
	assert.Equal(t, svc, model.DynamoDBAPI, "DynamoDB service from model is not matching the one used as parameter")
}
