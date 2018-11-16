package model

import (
	"errors"
	"fmt"
	"regexp"

	"gitlab.com/auto-staging/builder/helper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/types"
	yaml "gopkg.in/yaml.v2"
)

func CreateCodeBuildJob(event types.Event) error {
	envVars := []*codebuild.EnvironmentVariable{}
	for key, value := range event.EnvironmentVariables {
		envVars = append(envVars, &codebuild.EnvironmentVariable{
			Name:  aws.String(key),
			Type:  aws.String("PLAINTEXT"),
			Value: aws.String(value),
		})
	}

	buildspec := types.Buildspec{
		Version: "0.2",
		Phases: types.Phases{
			Build: types.Build{
				Commands: []string{
					"terraform --version",
				},
				Finally: []string{
					"aws lambda invoke --function-name auto-staging-builder --invocation-type Event --payload '{ \"operation\": \"RESULT_CREATE\", \"success\": '${CODEBUILD_BUILD_SUCCEEDING}', \"repository\": \"" + event.Repository + "\", \"branch\": \"" + event.Branch + "\" }'  /dev/null",
				},
			},
		},
	}

	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "regex/compile"}, 0)
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	res, err := yaml.Marshal(buildspec)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "yaml/marshal"}, 0)
		return err
	}

	helper.Logger.Log(errors.New(fmt.Sprint(string(res))), map[string]string{"module": "model/CreateCodeBuildJob", "operation": "buildspec"}, 4)

	createInput := codebuild.CreateProjectInput{
		Name:        aws.String("auto-staging-" + event.Repository + "-" + branchName),
		Description: aws.String("Managed by auto-staging"),
		ServiceRole: aws.String("arn:aws:iam::171842373341:role/auto-staging-builder-codebuild-exec-role"),
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          aws.String("BUILD_GENERAL1_SMALL"),
			Image:                aws.String("janrtr/auto-staging-build"),
			Type:                 aws.String("LINUX_CONTAINER"),
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      aws.String("GITHUB"),
			Location:  aws.String(event.RepositoryURL),
			Buildspec: aws.String(string(res)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String("NO_ARTIFACTS"),
		},
	}

	client := getCodeBuildClient()
	_, err = client.CreateProject(&createInput)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "codebuild/create"}, 0)
		return err
	}

	return err
}

func SetStatusAfterCreation(event types.Event) error {

	status := "init failed"

	if event.Success == 1 {
		status = "running"
	}

	return setStatusForEnvironment(event, status)
}
