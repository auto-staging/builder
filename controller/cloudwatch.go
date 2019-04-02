package controller

import (
	"os"

	"github.com/auto-staging/builder/helper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func getCloudWatchEventsClient() *cloudwatchevents.CloudWatchEvents {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/getCloudWatchClient", "operation": "aws/session"}, 0)
	}

	return cloudwatchevents.New(sess)
}
