package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
)

type CloudWatchEventsModelAPI interface {
	UpdateCloudWatchEvents(event types.Event) error
	DeleteCloudWatchEvents(event types.Event) error
}

type CloudWatchEventsModel struct {
	cloudwatcheventsiface.CloudWatchEventsAPI
}

func NewCloudWatchEventsModel(svc cloudwatcheventsiface.CloudWatchEventsAPI) *CloudWatchEventsModel {
	return &CloudWatchEventsModel{
		CloudWatchEventsAPI: svc,
	}
}
