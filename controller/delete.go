package controller

import (
	"errors"
	"fmt"

	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteController(event types.Event) (string, error) {

	status := types.Status{}
	err := model.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "running" && status.Status != "stopped" && status.Status != "initiating failed" {
		helper.Logger.Log(errors.New("Can't delete environment in status = "+status.Status), map[string]string{"module": "controller/DeleteController", "operation": "statusCheck"}, 0)
		return fmt.Sprintf(fmt.Sprint("{\"message\" : \"can't delete environment in current status\"}")), err
	}

	err = model.AdaptCodeBildJobForDelete(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	err = model.TriggerCodeBuild(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}

func DeleteCloudWatchEventController(event types.Event) (string, error) {

	err := model.DeleteCloudWatchEvents(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}

func DeleteResultController(event types.Event) (string, error) {

	err := model.SetStatusAfterDeletion(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if event.Success == 1 {
		err = model.DeleteCodeBuildJob(event)
		if err != nil {
			return fmt.Sprintf(""), err
		}
		err = model.DeleteEnvironment(event)
		if err != nil {
			return fmt.Sprintf(""), err
		}
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}
