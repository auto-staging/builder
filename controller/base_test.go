package controller

import (
	"testing"

	"github.com/auto-staging/builder/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewServiceBaseController(t *testing.T) {
	modelCloudWatchEvents := new(mocks.CloudWatchEventsModelAPI)
	modelCodeBuild := new(mocks.CodeBuildModelAPI)
	modelDynamoDB := new(mocks.DynamoDBModelAPI)

	controller := NewServiceBaseController(modelCloudWatchEvents, modelCodeBuild, modelDynamoDB)

	assert.NotEmpty(t, controller, "Expected controller not to be empty")
	assert.Equal(t, modelCloudWatchEvents, controller.CloudWatchEventsModelAPI, "CloudWatchEventsModelAPI model from controller is not matching the one used as parameter")
	assert.Equal(t, modelCodeBuild, controller.CodeBuildModelAPI, "CodeBuildModelAPI model from controller is not matching the one used as parameter")
	assert.Equal(t, modelDynamoDB, controller.DynamoDBModelAPI, "DynamoDBModelAPI model from controller is not matching the one used as parameter")
}
