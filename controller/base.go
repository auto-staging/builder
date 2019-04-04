package controller

import (
	"github.com/auto-staging/builder/model"
	"github.com/auto-staging/builder/types"
)

// ServiceBaseControllerAPI is an interface including all controller functions
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

// ServiceBaseController is a struct including the builder model interfaces, all controller functions are called on this struct and the included models
type ServiceBaseController struct {
	model.CloudWatchEventsModelAPI
	model.CodeBuildModelAPI
	model.DynamoDBModelAPI
}

// NewServiceBaseController takes the builder model interfaces as parameter and returns the pointer to an ServiceBaseController struct, on which all controller functions with their model calls can be called
func NewServiceBaseController(modelCloudWatchEvents model.CloudWatchEventsModelAPI, modelCodeBuild model.CodeBuildModelAPI, modelDynamoDB model.DynamoDBModelAPI) *ServiceBaseController {
	return &ServiceBaseController{
		CloudWatchEventsModelAPI: modelCloudWatchEvents,
		CodeBuildModelAPI:        modelCodeBuild,
		DynamoDBModelAPI:         modelDynamoDB,
	}
}
