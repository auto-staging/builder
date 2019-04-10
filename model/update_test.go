package model

import (
	"testing"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/mocks"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdaptCodeBildJobForUpdate(t *testing.T) {
	helper.Init()
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
		EnvironmentVariables: []types.EnvironmentVariable{
			types.EnvironmentVariable{
				Name:  "testVar",
				Type:  "PLAINTEXT",
				Value: "testVarValue",
			},
		},
	}

	checkParameters := func(input *codebuild.UpdateProjectInput) error {
		environmentVariables := input.Environment.EnvironmentVariables
		foundBranchRaw := false
		foundBranch := false
		foundRepository := false
		foundRandom := false
		foundTestVar := false
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
				if *value.Value != "random123" {
					t.Error("TF_VAR_random should have been random123, was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "testVar":
				foundTestVar = true
				if *value.Value != "testVarValue" {
					t.Error("testVar should have been testVarValue, was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			}
		}
		if !foundBranch || !foundBranchRaw || !foundRepository || !foundRandom {
			t.Errorf("Expected all default environment variables to exist, was TF_VAR_branch_raw = %t, TF_VAR_branch = %t, TF_VAR_repository = %t , TF_VAR_random = %t, testVar = %t", foundBranchRaw, foundBranch, foundRepository, foundRandom, foundTestVar)
			t.FailNow()
		}
		return nil
	}

	svc := new(mocks.CodeBuildAPI)
	svc.On("UpdateProject", mock.AnythingOfType("*codebuild.UpdateProjectInput")).Return(nil, checkParameters)

	svc.On("BatchGetProjects", mock.AnythingOfType("*codebuild.BatchGetProjectsInput")).Return(&codebuild.BatchGetProjectsOutput{
		Projects: []*codebuild.Project{
			&codebuild.Project{
				Environment: &codebuild.ProjectEnvironment{
					EnvironmentVariables: []*codebuild.EnvironmentVariable{
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_repository"),
							Value: aws.String(event.Repository),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_branch"),
							Value: aws.String(event.Branch),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_random"),
							Value: aws.String("random123"),
						},
					},
				},
				Source: &codebuild.ProjectSource{
					Type: aws.String("GITHUB"),
				},
				Artifacts: &codebuild.ProjectArtifacts{
					Type: aws.String("NO_ARTIFACTS"),
				},
			},
		},
	}, nil)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := codeBuildModel.AdaptCodeBildJobForUpdate(event)

	assert.Nil(t, err, "Expected no error")
}

func TestAdaptCodeBildJobForUpdateGetError(t *testing.T) {
	helper.Init()

	event := types.Event{
		Branch:     "testBranch",
		Repository: "testRepo",
	}

	svc := new(mocks.CodeBuildAPI)
	awsError := errors.New("AWS SDK Test Error")
	svc.On("UpdateProject", mock.AnythingOfType("*codebuild.UpdateProjectInput")).Return(nil, nil)

	svc.On("BatchGetProjects", mock.AnythingOfType("*codebuild.BatchGetProjectsInput")).Return(&codebuild.BatchGetProjectsOutput{
		Projects: []*codebuild.Project{
			&codebuild.Project{
				Environment: &codebuild.ProjectEnvironment{
					EnvironmentVariables: []*codebuild.EnvironmentVariable{
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_repository"),
							Value: aws.String(event.Repository),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_branch"),
							Value: aws.String(event.Branch),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_random"),
							Value: aws.String("random123"),
						},
					},
				},
				Source: &codebuild.ProjectSource{
					Type: aws.String("GITHUB"),
				},
				Artifacts: &codebuild.ProjectArtifacts{
					Type: aws.String("NO_ARTIFACTS"),
				},
			},
		},
	}, awsError)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := codeBuildModel.AdaptCodeBildJobForUpdate(event)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, err, awsError)
}

func TestAdaptCodeBildJobForUpdateUpdateError(t *testing.T) {
	helper.Init()

	event := types.Event{
		Branch:     "testBranch",
		Repository: "testRepo",
	}

	svc := new(mocks.CodeBuildAPI)
	awsError := errors.New("AWS SDK Test Error")
	svc.On("UpdateProject", mock.AnythingOfType("*codebuild.UpdateProjectInput")).Return(nil, awsError)

	svc.On("BatchGetProjects", mock.AnythingOfType("*codebuild.BatchGetProjectsInput")).Return(&codebuild.BatchGetProjectsOutput{
		Projects: []*codebuild.Project{
			&codebuild.Project{
				Environment: &codebuild.ProjectEnvironment{
					EnvironmentVariables: []*codebuild.EnvironmentVariable{
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_repository"),
							Value: aws.String(event.Repository),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_branch"),
							Value: aws.String(event.Branch),
						},
						&codebuild.EnvironmentVariable{
							Name:  aws.String("TF_VAR_random"),
							Value: aws.String("random123"),
						},
					},
				},
				Source: &codebuild.ProjectSource{
					Type: aws.String("GITHUB"),
				},
				Artifacts: &codebuild.ProjectArtifacts{
					Type: aws.String("NO_ARTIFACTS"),
				},
			},
		},
	}, nil)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := codeBuildModel.AdaptCodeBildJobForUpdate(event)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, err, awsError)
}

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
