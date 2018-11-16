package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func setStatusForEnvironment(event types.Event, status string) error {
	svc := getDynamoDbClient()

	updateStruct := types.StatusUpdate{
		Status: status,
	}

	update, err := dynamodbattribute.MarshalMap(updateStruct)

	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/setStatusForEnvironment", "operation": "marshal"}, 0)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-environments"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(event.Repository),
			},
			"branch": {
				S: aws.String(event.Branch),
			},
		},
		UpdateExpression:          aws.String("SET #status = :status"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository) AND attribute_exists(branch)"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/setStatusForEnvironment", "operation": "dynamodb/exec"}, 0)
		return err
	}

	return err
}
