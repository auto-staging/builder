package model

import (
	"os"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// CreateModelAPI is an interface including all Create model functions
type DatabaseModelAPI interface {
	GetStatusForEnvironment(event types.Event, status *types.Status) error
}

type DatabaseModel struct {
	dynamodbiface.DynamoDBAPI
}

func NewDatabaseModel(svc dynamodbiface.DynamoDBAPI) *DatabaseModel {
	return &DatabaseModel{
		DynamoDBAPI: svc,
	}
}

func getDynamoDbClient() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/getDynamoDbClient", "operation": "aws/session"}, 0)
	}

	return dynamodb.New(sess)
}
