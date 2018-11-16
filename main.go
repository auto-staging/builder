package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/auto-staging/builder/controller"
	"gitlab.com/auto-staging/builder/helper"

	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/auto-staging/builder/types"
)

func HandleRequest(ctx context.Context, event types.Event) (string, error) {

	if event.Operation == "CREATE" {
		return controller.CreateController(event)
	}

	if event.Operation == "DELETE" {
		return controller.DeleteController(event)
	}

	if event.Operation == "RESULT_CREATE" {
		return controller.CreateResultController(event)
	}

	if event.Operation == "RESULT_DESTROY" {
		return controller.DeleteResultController(event)
	}

	log.Println("UNKNOWN OPERATION")
	return fmt.Sprintf("{\"message\": \"unknown operation\"}"), nil
}

func main() {
	helper.Init()

	lambda.Start(HandleRequest)
}
