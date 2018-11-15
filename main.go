package main

import (
	"context"
	"fmt"

	"gitlab.com/auto-staging/builder/controller"

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

	return fmt.Sprintf("{\"message\": \"unknown operation\"}"), nil
}

func main() {
	lambda.Start(HandleRequest)
}
