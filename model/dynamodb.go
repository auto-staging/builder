package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DatabaseModelAPI is an interface including all DynamoDB model functions
type DatabaseModelAPI interface {
	SetStatusAfterCreation(event types.Event) error
	SetStatusAfterDeletion(event types.Event) error
	DeleteEnvironment(event types.Event) error
	GetStatusForEnvironment(event types.Event, status *types.Status) error
	SetStatusForEnvironment(event types.Event, status string) error
	SetStatusAfterUpdate(event types.Event) error
}

type DatabaseModel struct {
	dynamodbiface.DynamoDBAPI
}

func NewDatabaseModel(svc dynamodbiface.DynamoDBAPI) *DatabaseModel {
	return &DatabaseModel{
		DynamoDBAPI: svc,
	}
}
