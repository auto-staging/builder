package model

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func TriggerCodeBuild(event types.Event) error {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/TriggerCodeBuild", "operation": "regex/compile"}, 0)
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	service := getCodeBuildClient()

	_, err = service.StartBuild(&codebuild.StartBuildInput{
		ProjectName: aws.String("auto-staging-" + event.Repository + "-" + branchName),
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/TriggerCodeBuild", "operation": "codebuild/exec"}, 0)
		return err
	}

	return err
}
