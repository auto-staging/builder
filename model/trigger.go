package model

import (
	"regexp"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
)

// TriggerCodeBuild starts the CodeBuild Job for the Environment specified in Event struct.
// If an error occurs the error gets logged and the returned.
func (CodeBuildModel *CodeBuildModel) TriggerCodeBuild(event types.Event) error {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/TriggerCodeBuild", "operation": "regex/compile"}, 0)
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	service := CodeBuildModel.CodeBuildAPI

	_, err = service.StartBuild(&codebuild.StartBuildInput{
		ProjectName: aws.String("auto-staging-" + event.Repository + "-" + branchName),
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/TriggerCodeBuild", "operation": "codebuild/exec"}, 0)
		return err
	}

	return err
}
