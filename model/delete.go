package model

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteCodeBuildJob(event types.Event) error {
	client := getCodeBuildClient()

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "controller/DeleteController", "operation": "regex/compile"}, 0)
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

func SetStatusAfterDeletion(event types.Event) error {

	status := "destroying failed"

	if event.Success == 1 {
		status = "destroyed"
	}

	return setStatusForEnvironment(event, status)
}
