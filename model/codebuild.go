package model

import (
	"github.com/aws/aws-sdk-go/service/codebuild/codebuildiface"
)

type CodeBuildModelAPI interface {
}

type CodeBuildModel struct {
	codebuildiface.CodeBuildAPI
}

func NewCodeBuildModel(svc codebuildiface.CodeBuildAPI) *CodeBuildModel {
	return &CodeBuildModel{
		CodeBuildAPI: svc,
	}
}
