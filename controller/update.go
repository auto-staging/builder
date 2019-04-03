package controller

import (
	"errors"
	"fmt"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/model"
	"github.com/auto-staging/builder/types"
)

// UpdateController is the controller for the UPDATE action.
// First the status of the Environment gets checked, if the status is "running" or "updating failed" the CodBuild Job gets adapted with the updated
// configuration and then triggered.
func (ServiceBaseController *ServiceBaseController) UpdateController(event types.Event) (string, error) {
	DynamoDBModel := model.NewDynamoDBModel(ServiceBaseController.DynamoDBAPI)
	codeBuildModel := model.NewCodeBuildModel(ServiceBaseController.CodeBuildAPI)

	status := types.Status{}
	err := DynamoDBModel.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "running" && status.Status != "updating failed" {
		helper.Logger.Log(errors.New("Can't update environment in status = "+status.Status), map[string]string{"module": "controller/UpdateController", "operation": "statusCheck"}, 0)
		return fmt.Sprint("{\"message\" : \"can't update environment in current status\"}"), err
	}

	err = DynamoDBModel.SetStatusForEnvironment(event, "updating")
	if err != nil {
		return "", err
	}
	err = codeBuildModel.AdaptCodeBildJobForUpdate(event)
	if err != nil {
		errStatus := DynamoDBModel.SetStatusForEnvironment(event, "updating failed")
		if errStatus != nil {
			return "", errStatus
		}
		return "", err
	}

	err = codeBuildModel.TriggerCodeBuild(event)
	if err != nil {
		return "", err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// UpdateResultController is the controller for the RESULT_UPDATE action.
// The status of the Environment gets set according to the result of the CodeBuild Job.
func (ServiceBaseController *ServiceBaseController) UpdateResultController(event types.Event) (string, error) {
	DynamoDBModel := model.NewDynamoDBModel(ServiceBaseController.DynamoDBAPI)

	err := DynamoDBModel.SetStatusAfterUpdate(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// UpdateCloudWatchEventController is the controller for the UPDATE_SCHEDULE action.
// It calls the function to update all CloudWatchEvents rules for the Environment.
func (ServiceBaseController *ServiceBaseController) UpdateCloudWatchEventController(event types.Event) (string, error) {
	cloudWatchEventsModel := model.NewCloudWatchEventsModel(ServiceBaseController.CloudWatchEventsAPI)

	err := cloudWatchEventsModel.UpdateCloudWatchEvents(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}
