package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
)

type CodeBuildModelAPI interface {
	CreateCodeBuildJob(event types.Event) error
	DeleteCodeBuildJob(event types.Event) error
	AdaptCodeBildJobForDelete(event types.Event) error
	TriggerCodeBuild(event types.Event) error
	AdaptCodeBildJobForUpdate(event types.Event) error
}

type CodeBuildModel struct {
	codebuildiface.CodeBuildAPI
}

func NewCodeBuildModel(svc codebuildiface.CodeBuildAPI) *CodeBuildModel {
	return &CodeBuildModel{
		CodeBuildAPI: svc,
	}
}
