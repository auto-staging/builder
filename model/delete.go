package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteCodeBuildJob(event types.Event) error {
	client := getCodeBuildClient()

	_, err := client.DeleteProject(&codebuild.DeleteProjectInput{
		Name: aws.String("auto-staging-" + event.Repository + "-" + event.Branch),
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/DeleteCodeBuildJob", "operation": "dynamodb/exec"}, 0)
	}

	return err
}
