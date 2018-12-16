package controller

import (
	"errors"
	"fmt"

	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func CreateController(event types.Event) (string, error) {

	status := types.Status{}
	err := model.GetStatusForEnvironment(event, &status)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	if status.Status != "pending" {
		helper.Logger.Log(errors.New("Can't create environment in status = "+status.Status), map[string]string{"module": "controller/CreateController", "operation": "statusCheck"}, 0)
		return fmt.Sprintf(fmt.Sprint("{\"message\" : \"can't create environment in current status\"}")), err
	}

	err = model.CreateCodeBuildJob(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	err = model.TriggerCodeBuild(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}

func CreateResultController(event types.Event) (string, error) {

	err := model.SetStatusAfterCreation(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}
