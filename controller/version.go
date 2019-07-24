package controller

import (
	"encoding/json"
	"fmt"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
)

// GetVersionController is the controller for the VERSION operation.
// It gets the version information which are stored in the binary on compilation and returns them as JSON
func GetVersionController(event types.Event) (string, error) {
	builderVersion := types.SingleComponentVersion{}
	helper.GetVersionInformation(&builderVersion)

	body, err := json.Marshal(builderVersion)
	if err != nil {
		helper.Logger.Log(err, map[string]string{"module": "controller/GetVersionController", "operation": "marshal"}, 0)
		return fmt.Sprint("{\"message\" : \"Internal server error\"}"), err
	}

	return string(body), nil
}
