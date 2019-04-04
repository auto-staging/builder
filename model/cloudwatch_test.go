package model

import (
	"testing"

	"github.com/auto-staging/builder/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudWatchEventsModel(t *testing.T) {
	svc := new(mocks.CloudWatchEventsAPI)

	model := NewCloudWatchEventsModel(svc)

	assert.NotEmpty(t, model, "Expected model not to be empty")
	assert.Equal(t, svc, model.CloudWatchEventsAPI, "CloudWatchEventsAPI service from model is not matching the one used as parameter")
}
