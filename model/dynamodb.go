package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DynamoDBModelAPI is an interface including all DynamoDB model functions
type DynamoDBModelAPI interface {
	SetStatusAfterCreation(event types.Event) error
	SetStatusAfterDeletion(event types.Event) error
	DeleteEnvironment(event types.Event) error
	GetStatusForEnvironment(event types.Event, status *types.Status) error
	SetStatusForEnvironment(event types.Event, status string) error
	SetStatusAfterUpdate(event types.Event) error
}

// DynamoDBModel is a struct including the AWS SDK DynamoDB interface, all DynamoDB model functions are called on this struct and the included AWS SDK DynamoDB service
type DynamoDBModel struct {
	dynamodbiface.DynamoDBAPI
}

// NewDynamoDBModel takes the AWS SDK DynamoDB interface as parameter and returns the pointer to an DynamoDBModel struct, on which all DynamoDB model functions can be called
func NewDynamoDBModel(svc dynamodbiface.DynamoDBAPI) *DynamoDBModel {
	return &DynamoDBModel{
		DynamoDBAPI: svc,
	}
}
