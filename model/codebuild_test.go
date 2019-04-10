package model

import (
	"testing"

	"github.com/auto-staging/builder/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewCodeBuildModel(t *testing.T) {
	svc := new(mocks.CodeBuildAPI)

	model := NewCodeBuildModel(svc)

	assert.NotEmpty(t, model, "Expected model not to be empty")
	assert.Equal(t, svc, model.CodeBuildAPI, "CodeBuildAPI service from model is not matching the one used as parameter")
}
