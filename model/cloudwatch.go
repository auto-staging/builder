package model

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
)

type CloudWatchEventsModelAPI interface {
}

type CloudWatchEventsModel struct {
	cloudwatcheventsiface.CloudWatchEventsAPI
}

func NewCloudWatchEventsModel(svc cloudwatcheventsiface.CloudWatchEventsAPI) *CloudWatchEventsModel {
	return &CloudWatchEventsModel{
		CloudWatchEventsAPI: svc,
	}
}
