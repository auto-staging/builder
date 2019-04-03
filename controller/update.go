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
func UpdateController(event types.Event) (string, error) {
	databaseModel := model.NewDatabaseModel(getDynamoDbClient())

	status := types.Status{}
	err := databaseModel.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "running" && status.Status != "updating failed" {
		helper.Logger.Log(errors.New("Can't update environment in status = "+status.Status), map[string]string{"module": "controller/UpdateController", "operation": "statusCheck"}, 0)
		return fmt.Sprint("{\"message\" : \"can't update environment in current status\"}"), err
	}

	err = databaseModel.SetStatusForEnvironment(event, "updating")
	if err != nil {
		return "", err
	}
	err = model.AdaptCodeBildJobForUpdate(event)
	if err != nil {
		errStatus := databaseModel.SetStatusForEnvironment(event, "updating failed")
		if errStatus != nil {
			return "", errStatus
		}
		return "", err
	}

	err = model.TriggerCodeBuild(event)
	if err != nil {
		return "", err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// UpdateResultController is the controller for the RESULT_UPDATE action.
// The status of the Environment gets set according to the result of the CodeBuild Job.
func UpdateResultController(event types.Event) (string, error) {
	databaseModel := model.NewDatabaseModel(getDynamoDbClient())

	err := databaseModel.SetStatusAfterUpdate(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// UpdateCloudWatchEventController is the controller for the UPDATE_SCHEDULE action.
// It calls the function to update all CloudWatchEvents rules for the Environment.
func UpdateCloudWatchEventController(event types.Event) (string, error) {
	cloudWatchEventsModel := model.NewCloudWatchEventsModel(getCloudWatchEventsClient())

	err := cloudWatchEventsModel.UpdateCloudWatchEvents(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}
