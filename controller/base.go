package controller

import (
	"github.com/auto-staging/builder/model"
	"github.com/auto-staging/builder/types"
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
	model.CloudWatchEventsModelAPI
	model.CodeBuildModelAPI
	model.DynamoDBModelAPI
}

func NewServiceBaseController(modelCloudWatchEvents model.CloudWatchEventsModelAPI, modelCodeBuild model.CodeBuildModelAPI, modelDynamoDB model.DynamoDBModelAPI) *ServiceBaseController {
	return &ServiceBaseController{
		CloudWatchEventsModelAPI: modelCloudWatchEvents,
		CodeBuildModelAPI:        modelCodeBuild,
		DynamoDBModelAPI:         modelDynamoDB,
	}
}
