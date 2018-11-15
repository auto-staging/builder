package model

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteCodeBuildJob(event types.Event) error {
	client := getCodeBuildClient()

	_, err := client.DeleteProject(&codebuild.DeleteProjectInput{
		Name: aws.String("auto-staging-" + event.Repository + "-" + event.Branch),
	})
	if err != nil {
		log.Println(err)
	}

	return err
}
