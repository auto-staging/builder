package main

import (
	"context"
	"fmt"

	"gitlab.com/auto-staging/builder/controller"
	"gitlab.com/auto-staging/builder/helper"

	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/auto-staging/builder/types"
)

func HandleRequest(ctx context.Context, event types.Event) (string, error) {

	switch event.Operation {
	case "CREATE":
		return controller.CreateController(event)

	case "DELETE":
		return controller.DeleteController(event)

	case "UPDATE":
		return controller.UpdateController(event)

	case "RESULT_CREATE":
		return controller.CreateResultController(event)

	case "RESULT_DESTROY":
		return controller.DeleteResultController(event)

	case "RESULT_UPDATE":
		return controller.UpdateResultController(event)

	default:
		return fmt.Sprintf("{\"message\": \"unknown operation\"}"), nil
	}
}

func main() {
	helper.Init()

	lambda.Start(HandleRequest)
}
