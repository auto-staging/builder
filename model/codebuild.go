package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"gitlab.com/auto-staging/builder/helper"
)

func getCodeBuildClient() *codebuild.CodeBuild {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/getCodeBuildClient", "operation": "aws/session"}, 0)
	}

	return codebuild.New(sess)
}
