package model

import (
	"testing"

	"github.com/auto-staging/builder/mocks"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetStatusAfterUpdateSuccess(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
		Success:               1,
	}

	checkParameters := func(input *dynamodb.UpdateItemInput) error {
		expressionAttributeValues := input.ExpressionAttributeValues
		updateStruct := types.StatusUpdate{}

		dynamodbattribute.UnmarshalMap(expressionAttributeValues, &updateStruct)

		if updateStruct.Status != "running" {
			t.Error("Exptected new status value to be running, was " + updateStruct.Status)
			t.FailNow()
			return errors.New("")
		}
		return nil
	}

	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, checkParameters)

	dynamoDBModel := DynamoDBModel{
		DynamoDBAPI: svc,
	}

	err := dynamoDBModel.SetStatusAfterUpdate(event)

	assert.Nil(t, err, "Expected no error")
}

func TestSetStatusAfterUpdateFailed(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
		Success:               0,
	}

	checkParameters := func(input *dynamodb.UpdateItemInput) error {
		expressionAttributeValues := input.ExpressionAttributeValues
		updateStruct := types.StatusUpdate{}

		dynamodbattribute.UnmarshalMap(expressionAttributeValues, &updateStruct)

		if updateStruct.Status != "updating failed" {
			t.Error("Exptected new status value to be updating failed, was " + updateStruct.Status)
			t.FailNow()
			return errors.New("")
		}
		return nil
	}

	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, checkParameters)

	dynamoDBModel := DynamoDBModel{
		DynamoDBAPI: svc,
	}

	err := dynamoDBModel.SetStatusAfterUpdate(event)

	assert.Nil(t, err, "Expected no error")
}
