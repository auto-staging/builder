package model

import (
	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// GetStatusForEnvironment gets the status for an Environment specified in the Event struct from DynamoDB.
// If an error occurs the error gets logged and the returned.
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

	err = dynamodbattribute.UnmarshalMap(result.Item, status)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/GetStatusForEnvironment", "operation": "dynamodb/unmarshalMap"}, 0)
		return err
	}

	return nil
}
