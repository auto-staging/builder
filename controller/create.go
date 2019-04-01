package controller

import (
	"errors"
	"fmt"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/model"
	"github.com/auto-staging/builder/types"
)

// CreateController is the controller for the CREATE action.
// First the status of the Environment gets checked, if the status is "pending" the CodBuild Job gets created and then triggered.
func CreateController(event types.Event) (string, error) {
	databaseModel := model.NewDatabaseModel(getDynamoDbClient())

	status := types.Status{}
	err := databaseModel.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "pending" {
		helper.Logger.Log(errors.New("Can't create environment in status = "+status.Status), map[string]string{"module": "controller/CreateController", "operation": "statusCheck"}, 0)
		return fmt.Sprint("{\"message\" : \"can't create environment in current status\"}"), err
	}

	err = databaseModel.SetStatusForEnvironment(event, "initiating")
	if err != nil {
		return "", err
	}

	err = model.CreateCodeBuildJob(event)
	if err != nil {
		errStatus := databaseModel.SetStatusForEnvironment(event, "initiating failed")
		if errStatus != nil {
			return "", errStatus
		}
	}
	return fmt.Sprintf(""), err

	err = model.TriggerCodeBuild(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}

// CreateResultController is the controller for the RESULT_CREATE action.
// The status of the Environment gets set according to the result of the CodeBuild Job.
func CreateResultController(event types.Event) (string, error) {
	databaseModel := model.NewDatabaseModel(getDynamoDbClient())

	err := databaseModel.SetStatusAfterCreation(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprint("{\"message\" : \"success\"}"), err
}
