package controller

import (
	"errors"
	"fmt"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
)

// DeleteController is the controller function for the DELETE action.
// First the status of the Environment gets checked, if the status is "running", "stopped", "initiating failed", "destroyed failed"
// the CodBuild Job gets adapted to delete the Environment and then triggered.
func (ServiceBaseController *ServiceBaseController) DeleteController(event types.Event) (string, error) {
	dynamoDBModel := ServiceBaseController.DynamoDBModelAPI
	codeBuildModel := ServiceBaseController.CodeBuildModelAPI

	status := types.Status{}
	err := dynamoDBModel.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "running" && status.Status != "stopped" && status.Status != "initiating failed" && status.Status != "destroying failed" {
		helper.Logger.Log(errors.New("Can't delete environment in status = "+status.Status), map[string]string{"module": "controller/DeleteController", "operation": "statusCheck"}, 0)
		return fmt.Sprint("{\"message\" : \"can't delete environment in current status\"}"), err
	}

	err = dynamoDBModel.SetStatusForEnvironment(event, "destroying")
	if err != nil {
		errStatus := dynamoDBModel.SetStatusForEnvironment(event, "destroying failed")
		if errStatus != nil {
			return "", errStatus
		}
		return "", err
	}

	err = codeBuildModel.AdaptCodeBildJobForDelete(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	err = codeBuildModel.TriggerCodeBuild(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// DeleteCloudWatchEventController is the controller function for the DELETE_SCHEDULE action.
// It calls the function to delete all CloudWatchEvents rules for the Environment.
func (ServiceBaseController *ServiceBaseController) DeleteCloudWatchEventController(event types.Event) (string, error) {
	cloudWatchEventsModel := ServiceBaseController.CloudWatchEventsModelAPI

	err := cloudWatchEventsModel.DeleteCloudWatchEvents(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// DeleteResultController is the controller function for the RESULT_DESTROY action.
// The status of the Environment gets set according to the result of the CodeBuild Job and the CodeBuild Job and Environment get removed.
func (ServiceBaseController *ServiceBaseController) DeleteResultController(event types.Event) (string, error) {
	dynamoDBModel := ServiceBaseController.DynamoDBModelAPI
	codeBuildModel := ServiceBaseController.CodeBuildModelAPI

	err := dynamoDBModel.SetStatusAfterDeletion(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if event.Success == 1 {
		err = codeBuildModel.DeleteCodeBuildJob(event)
		if err != nil {
			return fmt.Sprintf(""), err
		}
		err = dynamoDBModel.DeleteEnvironment(event)
		if err != nil {
			return fmt.Sprintf(""), err
		}
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}
