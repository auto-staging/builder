package model

import (
	"testing"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/mocks"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCodeBuildJob(t *testing.T) {
	helper.Init()
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
	}

	checkParameters := func(input *codebuild.CreateProjectInput) error {
		environmentVariables := input.Environment.EnvironmentVariables
		foundBranchRaw := false
		foundBranch := false
		foundRepository := false
		foundRandom := false
		for _, value := range environmentVariables {
			switch *value.Name {
			case "TF_VAR_branch_raw":
				foundBranchRaw = true
				if *value.Value != event.Branch {
					t.Error("TF_VAR_branch_raw should have been " + event.Branch + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_branch":
				foundBranch = true
				if *value.Value != event.Branch {
					t.Error("TF_VAR_branch should have been " + event.Branch + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_repository":
				foundRepository = true
				if *value.Value != event.Repository {
					t.Error("TF_VAR_repository should have been " + event.Repository + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_random":
				foundRandom = true
			}
		}
		if !foundBranch || !foundBranchRaw || !foundRepository || !foundRandom {
			t.Errorf("Expected all default environment variables to exist, was TF_VAR_branch_raw = %t, TF_VAR_branch = %t, TF_VAR_repository = %t , TF_VAR_random = %t", foundBranchRaw, foundBranch, foundRepository, foundRandom)
			t.FailNow()
		}
		return nil
	}

	svc := new(mocks.CodeBuildAPI)
	svc.On("CreateProject", mock.AnythingOfType("*codebuild.CreateProjectInput")).Return(nil, checkParameters)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := codeBuildModel.CreateCodeBuildJob(event)

	assert.Nil(t, err, "Expected no error")
}

func TestCreateCodeBuildJobAWSError(t *testing.T) {
	helper.Init()

	svc := new(mocks.CodeBuildAPI)
	awsError := errors.New("AWS SDK Test Error")
	svc.On("CreateProject", mock.AnythingOfType("*codebuild.CreateProjectInput")).Return(nil, awsError)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	event := types.Event{
		Branch:     "testBranch",
		Repository: "testRepo",
	}

	err := codeBuildModel.CreateCodeBuildJob(event)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, err, awsError)
}

func TestSetStatusAfterCreationSuccess(t *testing.T) {
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

	err := dynamoDBModel.SetStatusAfterCreation(event)

	assert.Nil(t, err, "Expected no error")
}

func TestSetStatusAfterCreationFailed(t *testing.T) {
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

		if updateStruct.Status != "initiating failed" {
			t.Error("Exptected new status value to be initiating failed, was " + updateStruct.Status)
			t.FailNow()
			return errors.New("")
		}
		return nil
	}

	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, checkParameters)

	dynamoDBModel := DynamoDBModel{
		DynamoDBAPI: svc,
	}

	err := dynamoDBModel.SetStatusAfterCreation(event)

	assert.Nil(t, err, "Expected no error")
}
