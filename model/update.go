package model

import (
	"fmt"
	"regexp"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// SetStatusForEnvironment updates the status of an environment (which gets identified by branch and repository from the event
// parameter) in DynamoDB with the status given as parameter
func (DynamoDBModel *DynamoDBModel) SetStatusForEnvironment(event types.Event, status string) error {
	svc := DynamoDBModel.DynamoDBAPI

	updateStruct := types.StatusUpdate{
		Status: status,
	}

	update, err := dynamodbattribute.MarshalMap(updateStruct)

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/setStatusForEnvironment", "operation": "marshal"}, 0)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-environments"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(event.Repository),
			},
			"branch": {
				S: aws.String(event.Branch),
			},
		},
		UpdateExpression:          aws.String("SET #status = :status"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository) AND attribute_exists(branch)"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/setStatusForEnvironment", "operation": "dynamodb/exec"}, 0)
		return err
	}

	return err
}

// AdaptCodeBildJobForUpdate adapts the CodeBuild Job and buildspec with the updated Environment configuration (EnvironmentVariables).
// If an error occurs the error gets logged and the returned.
func (CodeBuildModel *CodeBuildModel) AdaptCodeBildJobForUpdate(event types.Event) error {
	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "regex/compile"}, 0)
		return err
	}

	var envVars []*codebuild.EnvironmentVariable
	// Set default variables
	envVars = append(envVars, &codebuild.EnvironmentVariable{
		Name:  aws.String("TF_VAR_branch_raw"),
		Type:  aws.String("PLAINTEXT"),
		Value: aws.String(event.Branch),
	})
	branchName := reg.ReplaceAllString(event.Branch, "-")
	envVars = append(envVars, &codebuild.EnvironmentVariable{
		Name:  aws.String("TF_VAR_branch"),
		Type:  aws.String("PLAINTEXT"),
		Value: aws.String(branchName),
	})
	envVars = append(envVars, &codebuild.EnvironmentVariable{
		Name:  aws.String("TF_VAR_repository"),
		Type:  aws.String("PLAINTEXT"),
		Value: aws.String(event.Repository),
	})

	for _, environmentVariable := range event.EnvironmentVariables {
		envVars = append(envVars, &codebuild.EnvironmentVariable{
			Name:  aws.String(environmentVariable.Name),
			Type:  aws.String(environmentVariable.Type),
			Value: aws.String(environmentVariable.Value),
		})
	}

	client := CodeBuildModel.CodeBuildAPI
	oldProjects, err := client.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{
			aws.String("auto-staging-" + event.Repository + "-" + branchName),
		},
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "codebuild/batchGetProjects"}, 0)
		return err
	}
	oldProject := oldProjects.Projects[0]

	// Reuse old TF_VAR_random value
	for _, oldEnvironmentVar := range oldProject.Environment.EnvironmentVariables {
		if *oldEnvironmentVar.Name == "TF_VAR_random" {
			envVars = append(envVars, oldEnvironmentVar)
		}
	}

	buildspec := types.Buildspec{
		Version: "0.2",
		Phases: types.Phases{
			Build: types.Build{
				Commands: []string{
					"make auto-staging-init",
					"make auto-staging-apply",
				},
				Finally: []string{
					"aws lambda invoke --function-name auto-staging-builder --invocation-type Event --payload '{ \"operation\": \"RESULT_UPDATE\", \"success\": '${CODEBUILD_BUILD_SUCCEEDING}', \"repository\": \"" + event.Repository + "\", \"branch\": \"" + event.Branch + "\" }'  /dev/null",
				},
			},
		},
	}
	marshaledBuildspec, err := yaml.Marshal(buildspec)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "yaml/marshal"}, 0)
		return err
	}
	helper.Logger.Log(errors.New(fmt.Sprint(string(marshaledBuildspec))), map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "buildspec"}, 4)

	_, err = client.UpdateProject(&codebuild.UpdateProjectInput{
		Name:        oldProject.Name,
		Description: oldProject.Description,
		ServiceRole: aws.String(event.CodeBuildRoleARN),
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          oldProject.Environment.ComputeType,
			Image:                oldProject.Environment.Image,
			Type:                 oldProject.Environment.Type,
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      oldProject.Source.Type,
			Location:  aws.String(event.InfrastructureRepoURL),
			Buildspec: aws.String(string(marshaledBuildspec)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: oldProject.Artifacts.Type,
		},
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "codebuild/update"}, 0)
		return err
	}

	return err
}

// SetStatusAfterUpdate checks the success variable in the event struct, which gets set in the CodeBuild Job. If success euqals 1 then the status
// gets set to "running" otherwise it gets set to "updating failed".
// If an error occurs the error gets logged and the returned.
func (DynamoDBModel *DynamoDBModel) SetStatusAfterUpdate(event types.Event) error {

	status := "updating failed"

	if event.Success == 1 {
		status = "running"
	}

	return DynamoDBModel.SetStatusForEnvironment(event, status)
}
