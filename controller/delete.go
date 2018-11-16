package controller

import (
	"fmt"

	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteController(event types.Event) (string, error) {

	err := model.AdaptCodeBildJobForDelete(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	err = model.TriggerCodeBuild(event)
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
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}
