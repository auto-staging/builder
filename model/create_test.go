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

func TestCreateCodeBuildJob(t *testing.T) {
	helper.Init()
	event := types.Event{
		Branch:                "testBranch",
		Repository:            "testRepo",
		InfrastructureRepoURL: "testUrl",
	}

	checkParameters := func(input *codebuild.CreateProjectInput) error {
		environmentVariables := input.Environment.EnvironmentVariables
		foundBranchRaw := false
		foundBranch := false
		foundRepository := false
		foundRandom := false
		for _, value := range environmentVariables {
			switch *value.Name {
			case "TF_VAR_branch_raw":
				foundBranchRaw = true
				if *value.Value != event.Branch {
					t.Error("TF_VAR_branch_raw should have been " + event.Branch + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_branch":
				foundBranch = true
				if *value.Value != event.Branch {
					t.Error("TF_VAR_branch should have been " + event.Branch + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_repository":
				foundRepository = true
				if *value.Value != event.Repository {
					t.Error("TF_VAR_repository should have been " + event.Repository + ", was " + *value.Value)
					t.FailNow()
					return errors.New("")
				}
			case "TF_VAR_random":
				foundRandom = true
			}
		}
		if !foundBranch || !foundBranchRaw || !foundRepository || !foundRandom {
			t.Errorf("Expected all default environment variables to exist, was TF_VAR_branch_raw = %t, TF_VAR_branch = %t, TF_VAR_repository = %t , TF_VAR_random = %t", foundBranchRaw, foundBranch, foundRepository, foundRandom)
			t.FailNow()
		}
		return nil
	}

	svc := new(mocks.CodeBuildAPI)
	svc.On("CreateProject", mock.AnythingOfType("*codebuild.CreateProjectInput")).Return(nil, checkParameters)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	err := codeBuildModel.CreateCodeBuildJob(event)

	assert.Nil(t, err, "Expected no error")
}

func TestCreateCodeBuildJobAWSError(t *testing.T) {
	helper.Init()

	svc := new(mocks.CodeBuildAPI)
	awsError := errors.New("AWS SDK Test Error")
	svc.On("CreateProject", mock.AnythingOfType("*codebuild.CreateProjectInput")).Return(nil, awsError)

	codeBuildModel := CodeBuildModel{
		CodeBuildAPI: svc,
	}

	event := types.Event{
		Branch:     "testBranch",
		Repository: "testRepo",
	}

	err := codeBuildModel.CreateCodeBuildJob(event)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, err, awsError)
}
