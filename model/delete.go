package model

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
	yaml "gopkg.in/yaml.v2"
)

func DeleteCodeBuildJob(event types.Event) error {
	client := getCodeBuildClient()

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

func AdaptCodeBildJobForDelete(event types.Event) error {
	err := setStatusForEnvironment(event, "destroying")
	if err != nil {
		return err
	}

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

	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "regex/compile"}, 0)
		setStatusForEnvironment(event, "destroying failed")
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	res, err := yaml.Marshal(buildspec)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "yaml/marshal"}, 0)
		setStatusForEnvironment(event, "destroying failed")
		return err
	}

	helper.Logger.Log(errors.New(fmt.Sprint(string(res))), map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "buildspec"}, 4)

	client := getCodeBuildClient()

	oldProjects, err := client.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{
			aws.String("auto-staging-" + event.Repository + "-" + branchName),
		},
	})

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "aws/batchGetProjects"}, 0)
		setStatusForEnvironment(event, "destroying failed")
		return err
	}

	oldProject := oldProjects.Projects[0]
	oldProject.Source.Buildspec = aws.String(string(res))

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
			Buildspec: aws.String(string(res)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: oldProject.Artifacts.Type,
		},
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptingCodeBildJobForDelete", "operation": "codebuild/update"}, 0)
		setStatusForEnvironment(event, "destroying failed")
		return err
	}

	return err
}

func SetStatusAfterDeletion(event types.Event) error {

	status := "destroying failed"

	if event.Success == 1 {
		status = "destroyed"
	}

	return setStatusForEnvironment(event, status)
}

func DeleteEnvironment(event types.Event) error {
	svc := getDynamoDbClient()

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
