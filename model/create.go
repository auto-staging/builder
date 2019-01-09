package model

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/auto-staging/builder/helper"

	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	yaml "gopkg.in/yaml.v2"
)

// CreateCodeBuildJob creates the CodebuildJob via AWS SDK with the configuration defined for the Environment.
// If an error occurs the error gets logged and the returned.
func CreateCodeBuildJob(event types.Event) error {
	err := setStatusForEnvironment(event, "initiating")
	if err != nil {
		return err
	}

	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "regex/compile"}, 0)
		errStatus := setStatusForEnvironment(event, "initiating failed")
		if errStatus != nil {
			return errStatus
		}
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	var envVars []*codebuild.EnvironmentVariable
	// Set default variables
	envVars = append(envVars, &codebuild.EnvironmentVariable{
		Name:  aws.String("TF_VAR_branch_raw"),
		Type:  aws.String("PLAINTEXT"),
		Value: aws.String(event.Branch),
	})
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

	buildspec := types.Buildspec{
		Version: "0.2",
		Phases: types.Phases{
			Build: types.Build{
				Commands: []string{
					"make auto-staging-init",
					"make auto-staging-apply",
				},
				Finally: []string{
					"aws lambda invoke --function-name auto-staging-builder --invocation-type Event --payload '{ \"operation\": \"RESULT_CREATE\", \"success\": '${CODEBUILD_BUILD_SUCCEEDING}', \"repository\": \"" + event.Repository + "\", \"branch\": \"" + event.Branch + "\" }'  /dev/null",
				},
			},
		},
	}

	res, err := yaml.Marshal(buildspec)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/CreateCodeBuildJob", "operation": "yaml/marshal"}, 0)
		errStatus := setStatusForEnvironment(event, "initiating failed")
		if errStatus != nil {
			return errStatus
		}
		return err
	}

	helper.Logger.Log(errors.New(fmt.Sprint(string(res))), map[string]string{"module": "model/CreateCodeBuildJob", "operation": "buildspec"}, 4)

	createInput := codebuild.CreateProjectInput{
		Name:        aws.String("auto-staging-" + event.Repository + "-" + branchName),
		Description: aws.String("Managed by auto-staging"),
		ServiceRole: aws.String(event.CodeBuildRoleARN),
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          aws.String("BUILD_GENERAL1_SMALL"),
			Image:                aws.String("janrtr/auto-staging-build"),
			Type:                 aws.String("LINUX_CONTAINER"),
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      aws.String("GITHUB"),
			Location:  aws.String(event.InfrastructureRepoURL),
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
		errStatus := setStatusForEnvironment(event, "initiating failed")
		if errStatus != nil {
			return errStatus
		}
		return err
	}

	return err
}

// SetStatusAfterCreation checks the success variable in the event struct, which gets set in the CodeBuild Job. If success euqals 1 then the status
// gets set to "running" otherwise it gets set to "initating failed".
// If an error occurs the error gets logged and the returned.
func SetStatusAfterCreation(event types.Event) error {

	status := "initiating failed"

	if event.Success == 1 {
		status = "running"
	}

	return setStatusForEnvironment(event, status)
}
