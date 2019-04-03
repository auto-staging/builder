package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
)

// CloudWatchEventsModelAPI is an interface including all CloudWatchEvents model functions
type CloudWatchEventsModelAPI interface {
	UpdateCloudWatchEvents(event types.Event) error
	DeleteCloudWatchEvents(event types.Event) error
}

// CloudWatchEventsModel is a struct including the AWS SDK CloudWatchEvents interface, all CloudWatchEvents model functions are called on this struct and the included AWS SDK CloudWatchEvents service
type CloudWatchEventsModel struct {
	cloudwatcheventsiface.CloudWatchEventsAPI
}

// NewCloudWatchEventsModel takes the AWS SDK CloudWatchEvents interface as parameter and returns the pointer to an CloudWatchEventsModel struct, on which all CloudWatchEvents model functions can be called
func NewCloudWatchEventsModel(svc cloudwatcheventsiface.CloudWatchEventsAPI) *CloudWatchEventsModel {
	return &CloudWatchEventsModel{
		CloudWatchEventsAPI: svc,
	}
}
