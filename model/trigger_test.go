package model

import (
	"testing"

	"github.com/auto-staging/builder/helper"
	"github.com/auto-staging/builder/mocks"
	"github.com/auto-staging/builder/types"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTriggerCodeBuild(t *testing.T) {
	svc := new(mocks.CodeBuildAPI)
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
		Success:               0,
	}

	expectedProjectName := "auto-staging-" + event.Repository + "-" + event.Branch

	checkParameters := func(input *codebuild.StartBuildInput) error {
		if *input.ProjectName != expectedProjectName {
			t.Error("Exptected project name to be " + expectedProjectName + ", was " + *input.ProjectName)
			t.FailNow()
			return errors.New("")
		}
		return nil
	}

	svc.On("StartBuild", mock.AnythingOfType("*codebuild.StartBuildInput")).Return(nil, checkParameters)

	CodeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := CodeBuildModel.TriggerCodeBuild(event)

	assert.Nil(t, err, "Expected no error")
}

func TestTriggerCodeBuildError(t *testing.T) {
	helper.Init()
	svc := new(mocks.CodeBuildAPI)
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
		Success:               0,
	}

	errAWS := errors.New("TestError")

	svc.On("StartBuild", mock.AnythingOfType("*codebuild.StartBuildInput")).Return(nil, errAWS)

	CodeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := CodeBuildModel.TriggerCodeBuild(event)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, err, errAWS, "Returned error didn't match the mocked error")
}
