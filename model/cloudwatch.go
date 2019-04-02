package model

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type CloudWatchModelAPI interface {
}

type CloudWatchModel struct {
	cloudwatchiface.CloudWatchAPI
}

func NewCloudWatchModel(svc cloudwatchiface.CloudWatchAPI) *CloudWatchModel {
	return &CloudWatchModel{
		CloudWatchAPI: svc,
	}
}
