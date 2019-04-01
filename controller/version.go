package controller

import (
	"encoding/json"
	"fmt"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/types"
)

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
