package controller

import (
	"fmt"
	"log"
	"regexp"

	"gitlab.com/auto-staging/builder/model"
	"gitlab.com/auto-staging/builder/types"
)

func CreateController(event types.Event) (string, error) {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	event.Branch = reg.ReplaceAllString(event.Branch, "-")

	err = model.CreateCodeBuildJob(event)
	if err != nil {
		return fmt.Sprintf(""), err
	}

	return fmt.Sprintf(fmt.Sprint("{\"message\" : \"success\"}")), err
}
