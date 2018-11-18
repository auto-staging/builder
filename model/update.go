package model

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
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

func AdaptCodeBildJobForUpdate(event types.Event) error {
	err := setStatusForEnvironment(event, "updating")
	if err != nil {
		return err
	}
	// Adapt branch name to only contain allowed characters for CodeBuild name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "regex/compile"}, 0)
		setStatusForEnvironment(event, "updating failed")
		return err
	}
	branchName := reg.ReplaceAllString(event.Branch, "-")

	client := getCodeBuildClient()
	oldProjects, err := client.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{
			aws.String("auto-staging-" + event.Repository + "-" + branchName),
		},
	})
	oldProject := oldProjects.Projects[0]

	envVars := []*codebuild.EnvironmentVariable{}
	for key, value := range event.EnvironmentVariables {
		envVars = append(envVars, &codebuild.EnvironmentVariable{
			Name:  aws.String(key),
			Type:  aws.String("PLAINTEXT"),
			Value: aws.String(value),
		})
	}

	_, err = client.UpdateProject(&codebuild.UpdateProjectInput{
		Name:        oldProject.Name,
		Description: oldProject.Description,
		ServiceRole: oldProject.ServiceRole,
		Environment: &codebuild.ProjectEnvironment{
			ComputeType:          oldProject.Environment.ComputeType,
			Image:                oldProject.Environment.Image,
			Type:                 oldProject.Environment.Type,
			EnvironmentVariables: envVars,
		},
		Source: &codebuild.ProjectSource{
			Type:      oldProject.Source.Type,
			Location:  aws.String(event.InfrastructureRepoUrl),
			Buildspec: oldProject.Source.Buildspec,
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: oldProject.Artifacts.Type,
		},
	})
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/AdaptCodeBildJobForUpdate", "operation": "codebuild/update"}, 0)
		setStatusForEnvironment(event, "updating failed")
		return err
	}

	return err
}

func SetStatusAfterUpdate(event types.Event) error {

	status := "updating failed"

	if event.Success == 1 {
		status = "running"
	}

	return setStatusForEnvironment(event, status)
}
