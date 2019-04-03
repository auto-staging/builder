package main

import (
	"context"
	"fmt"
	"os"

	"github.com/auto-staging/builder/controller"
	"github.com/auto-staging/builder/helper"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/auto-staging/builder/types"
)

// HandleRequest redirects the request to the matching controller based on the operation in the event.
func HandleRequest(ctx context.Context, event types.Event) (string, error) {

	serviceBaseController := controller.NewServiceBaseController(
		getCloudWatchEventsClient(),
		getCodeBuildClient(),
		getDynamoDbClient(),
	)

	switch event.Operation {
	case "CREATE":
		return serviceBaseController.CreateController(event)

	case "DELETE":
		return serviceBaseController.DeleteController(event)

	case "UPDATE":
		return serviceBaseController.UpdateController(event)

	case "RESULT_CREATE":
		return serviceBaseController.CreateResultController(event)

	case "RESULT_DESTROY":
		return serviceBaseController.DeleteResultController(event)

	case "RESULT_UPDATE":
		return serviceBaseController.UpdateResultController(event)

	case "UPDATE_SCHEDULE":
		return serviceBaseController.UpdateCloudWatchEventController(event)

	case "DELETE_SCHEDULE":
		return serviceBaseController.DeleteCloudWatchEventController(event)

	default:
		return fmt.Sprintf("{\"message\": \"unknown operation\"}"), nil
	}
}

func main() {
	helper.Init()

	lambda.Start(HandleRequest)
}

func getCloudWatchEventsClient() *cloudwatchevents.CloudWatchEvents {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "main/getCloudWatchClient", "operation": "aws/session"}, 0)
	}

	return cloudwatchevents.New(sess)
}

func getCodeBuildClient() *codebuild.CodeBuild {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "main/getCodeBuildClient", "operation": "aws/session"}, 0)
	}

	return codebuild.New(sess)
}

func getDynamoDbClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "main/getDynamoDbClient", "operation": "aws/session"}, 0)
	}

	return dynamodb.New(sess)
}
