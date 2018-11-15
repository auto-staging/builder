package controller

import (
	"fmt"
	"log"
	"regexp"

	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func DeleteController(event types.Event) (string, error) {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	event.Branch = reg.ReplaceAllString(event.Branch, "-")

	err = model.DeleteCodeBuildJob(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}
