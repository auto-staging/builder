package model

import (
	"fmt"
	"regexp"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// DeleteCodeBuildJob removes the CodeBuild Job for the Environment specified in the Event struct.
// If an error occurs the error gets logged and the returned.
func (CodeBuildModel *CodeBuildModel) DeleteCodeBuildJob(event types.Event) error {
	client := CodeBuildModel.CodeBuildAPI

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "controller/DeleteCodeBuildJob", "operation": "regex/compile"}, 0)
		return err
	}
	event.Branch = reg.ReplaceAllString(event.Branch, "-")

	_, err = client.DeleteProject(&codebuild.DeleteProjectInput{
		Name: aws.String("auto-staging-" + event.Repository + "-" + event.Branch),
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/DeleteCodeBuildJob", "operation": "dynamodb/exec"}, 0)
	}

	return err
}

// AdaptCodeBildJobForDelete adapts the CodeBuild buildspec to delete an Environment infrastructure.
// If an error occurs the error gets logged and the returned.
func (CodeBuildModel *CodeBuildModel) AdaptCodeBildJobForDelete(event types.Event) error {
	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "regex/compile"}, 0)
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	client := CodeBuildModel.CodeBuildAPI

	oldProjects, err := client.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{
			aws.String("auto-staging-" + event.Repository + "-" + branchName),
		},
	})

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "aws/batchGetProjects"}, 0)
		return err
	}

	oldProject := oldProjects.Projects[0]

	buildspec := types.Buildspec{
		Version: "0.2",
		Phases: types.Phases{
			Build: types.Build{
				Commands: []string{
					"make auto-staging-init",
					"make auto-staging-destroy",
				},
				Finally: []string{
					"aws lambda invoke --function-name auto-staging-builder --invocation-type Event --payload '{ \"operation\": \"RESULT_DESTROY\", \"success\": '${CODEBUILD_BUILD_SUCCEEDING}', \"repository\": \"" + event.Repository + "\", \"branch\": \"" + event.Branch + "\" }'  /dev/null",
				},
			},
		},
	}
	marshaledBuildspec, err := yaml.Marshal(buildspec)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "yaml/marshal"}, 0)
		return err
	}
	helper.Logger.Log(errors.New(fmt.Sprint(string(marshaledBuildspec))), map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "buildspec"}, 4)

	_, err = client.UpdateProject(&codebuild.UpdateProjectInput{
		Name:        oldProject.Name,
		Description: oldProject.Description,
		ServiceRole: oldProject.ServiceRole,
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          oldProject.Environment.ComputeType,
			Image:                oldProject.Environment.Image,
			Type:                 oldProject.Environment.Type,
			EnvironmentVariables: oldProject.Environment.EnvironmentVariables,
		},
		Source: &codebuild.ProjectSource{
			Type:      oldProject.Source.Type,
			Location:  oldProject.Source.Location,
			Buildspec: aws.String(string(marshaledBuildspec)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: oldProject.Artifacts.Type,
		},
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "codebuild/update"}, 0)
		return err
	}

	return err
}

// SetStatusAfterDeletion checks the success variable in the event struct, which gets set in the CodeBuild Job. If success euqals 1 then the status
// gets set to "destroyed" otherwise it gets set to "destroying failed".
// If an error occurs the error gets logged and the returned.
func (DynamoDBModel *DynamoDBModel) SetStatusAfterDeletion(event types.Event) error {

	status := "destroying failed"

	if event.Success == 1 {
		status = "destroyed"
	}

	return DynamoDBModel.SetStatusForEnvironment(event, status)
}

// DeleteEnvironment removes an Environment specified in the Event struct from DynamoDB.
// If an error occurs the error gets logged and the returned.
func (DynamoDBModel *DynamoDBModel) DeleteEnvironment(event types.Event) error {
	svc := DynamoDBModel.DynamoDBAPI

	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("auto-staging-environments"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(event.Repository),
			},
			"branch": {
				S: aws.String(event.Branch),
			},
		},
	})

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/DeleteEnvironment", "operation": "dynamodb/exec"}, 0)
		return err
	}

	return nil
}
