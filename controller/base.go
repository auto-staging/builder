package controller

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type ServiceBaseControllerAPI interface {
	CreateController(event types.Event) (string, error)
	CreateResultController(event types.Event) (string, error)
	DeleteController(event types.Event) (string, error)
	DeleteCloudWatchEventController(event types.Event) (string, error)
	DeleteResultController(event types.Event) (string, error)
	UpdateController(event types.Event) (string, error)
	UpdateResultController(event types.Event) (string, error)
	UpdateCloudWatchEventController(event types.Event) (string, error)
}

type ServiceBaseController struct {
	cloudwatcheventsiface.CloudWatchEventsAPI
	codebuildiface.CodeBuildAPI
	dynamodbiface.DynamoDBAPI
}

func NewServiceBaseController(svcCloudWatchEvents cloudwatcheventsiface.CloudWatchEventsAPI, svcCodeBuild codebuildiface.CodeBuildAPI, svcDynamoDB dynamodbiface.DynamoDBAPI) *ServiceBaseController {
	return &ServiceBaseController{
		CloudWatchEventsAPI: svcCloudWatchEvents,
		CodeBuildAPI:        svcCodeBuild,
		DynamoDBAPI:         svcDynamoDB,
	}
}
