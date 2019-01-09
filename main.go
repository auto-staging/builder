package main

import (
	"context"
	"fmt"

	"github.com/auto-staging/builder/controller"
	"github.com/auto-staging/builder/helper"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/auto-staging/builder/types"
)

// HandleRequest redirects the request to the matching controller based on the operation in the event.
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

	case "UPDATE_SCHEDULE":
		return controller.UpdateCloudWatchEventController(event)

	case "DELETE_SCHEDULE":
		return controller.DeleteCloudWatchEventController(event)

	default:
		return fmt.Sprintf("{\"message\": \"unknown operation\"}"), nil
	}
}

func main() {
	helper.Init()

	lambda.Start(HandleRequest)
}
