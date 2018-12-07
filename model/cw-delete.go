package model

import (
	"fmt"
	"regexp"

	"gitlab.com/auto-staging/builder/helper"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteCloudWatchEvents(event types.Event) error {
	fmt.Println(event)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "model/DeleteCloudWatchEvents", "operation": "regex/compile"}, 0)
		return err
	}

	branchName := reg.ReplaceAllString(event.Branch, "-")

	// Startup schedules
	err = removeRulesWithTarget(event.Repository, branchName, event.Branch, "start")
	if err != nil {
		return err
	}

	// Shutdown schedules
	err = removeRulesWithTarget(event.Repository, branchName, event.Branch, "stop")
	if err != nil {
		return err
	}

	return nil
}
