package model

import (
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
)

// CodeBuildModelAPI is an interface including all CodeBuild model functions
type CodeBuildModelAPI interface {
	CreateCodeBuildJob(event types.Event) error
	DeleteCodeBuildJob(event types.Event) error
	AdaptCodeBildJobForDelete(event types.Event) error
	TriggerCodeBuild(event types.Event) error
	AdaptCodeBildJobForUpdate(event types.Event) error
}

// CodeBuildModel is a struct including the AWS SDK CodeBuild interface, all CodeBuild model functions are called on this struct and the included AWS SDK CodeBuild service
type CodeBuildModel struct {
	codebuildiface.CodeBuildAPI
}

// NewCodeBuildModel takes the AWS SDK CodeBuild interface as parameter and returns the pointer to an CodeBuildModel struct, on which all CodeBuild model functions can be called
func NewCodeBuildModel(svc codebuildiface.CodeBuildAPI) *CodeBuildModel {
	return &CodeBuildModel{
		CodeBuildAPI: svc,
	}
}
