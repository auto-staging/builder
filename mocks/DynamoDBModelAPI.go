// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

import types "github.com/auto-staging/builder/types"

// DynamoDBModelAPI is an autogenerated mock type for the DynamoDBModelAPI type
type DynamoDBModelAPI struct {
	mock.Mock
}

// DeleteEnvironment provides a mock function with given fields: event
func (_m *DynamoDBModelAPI) DeleteEnvironment(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStatusForEnvironment provides a mock function with given fields: event, status
func (_m *DynamoDBModelAPI) GetStatusForEnvironment(event types.Event, status *types.Status) error {
	ret := _m.Called(event, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event, *types.Status) error); ok {
		r0 = rf(event, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStatusAfterCreation provides a mock function with given fields: event
func (_m *DynamoDBModelAPI) SetStatusAfterCreation(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStatusAfterDeletion provides a mock function with given fields: event
func (_m *DynamoDBModelAPI) SetStatusAfterDeletion(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStatusAfterUpdate provides a mock function with given fields: event
func (_m *DynamoDBModelAPI) SetStatusAfterUpdate(event types.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStatusForEnvironment provides a mock function with given fields: event, status
func (_m *DynamoDBModelAPI) SetStatusForEnvironment(event types.Event, status string) error {
	ret := _m.Called(event, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Event, string) error); ok {
		r0 = rf(event, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
