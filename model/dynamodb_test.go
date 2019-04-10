package model

import (
	"testing"

	"github.com/auto-staging/builder/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewDynamoDBModel(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)

	model := NewDynamoDBModel(svc)

	assert.NotEmpty(t, model, "Expected model not to be empty")
	assert.Equal(t, svc, model.DynamoDBAPI, "DynamoDB service from model is not matching the one used as parameter")
}
