package controller

import (
	"fmt"

	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func CreateController(event types.Event) (string, error) {

	err := model.CreateCodeBuildJob(event)
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
