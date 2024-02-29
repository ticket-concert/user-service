// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	request "user-service/internal/modules/user/models/request"

	mock "github.com/stretchr/testify/mock"

	response "user-service/internal/modules/user/models/response"
)

// UsecaseQuery is an autogenerated mock type for the UsecaseQuery type
type UsecaseQuery struct {
	mock.Mock
}

// GetProfile provides a mock function with given fields: origCtx, payload
func (_m *UsecaseQuery) GetProfile(origCtx context.Context, payload request.GetProfile) (*response.GetProfile, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for GetProfile")
	}

	var r0 *response.GetProfile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.GetProfile) (*response.GetProfile, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.GetProfile) *response.GetProfile); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.GetProfile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.GetProfile) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsecaseQuery creates a new instance of UsecaseQuery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecaseQuery(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsecaseQuery {
	mock := &UsecaseQuery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
