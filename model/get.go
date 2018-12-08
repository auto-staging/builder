package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func GetStatusForEnvironment(event types.Event, status *types.Status) error {
	svc := getDynamoDbClient()

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("auto-staging-environments"),
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(event.Repository),
			},
			"branch": {
				S: aws.String(event.Branch),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		ProjectionExpression: aws.String("#status"),
	})

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/GetStatusForEnvironment", "operation": "dynamodb/exec"}, 0)
		return err
	}

	dynamodbattribute.UnmarshalMap(result.Item, status)

	return nil
}
