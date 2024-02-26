// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"
	request "user-service/internal/modules/user/models/request"

	mock "github.com/stretchr/testify/mock"

	response "user-service/internal/modules/user/models/response"
)

// UsecaseCommand is an autogenerated mock type for the UsecaseCommand type
type UsecaseCommand struct {
	mock.Mock
}

// LoginUser provides a mock function with given fields: origCtx, payload
func (_m *UsecaseCommand) LoginUser(origCtx context.Context, payload request.LoginUser) (*response.LoginUserResp, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for LoginUser")
	}

	var r0 *response.LoginUserResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.LoginUser) (*response.LoginUserResp, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.LoginUser) *response.LoginUserResp); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.LoginUserResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.LoginUser) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: origCtx, payload
func (_m *UsecaseCommand) RegisterUser(origCtx context.Context, payload request.RegisterUser) (*response.RegisterUser, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 *response.RegisterUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.RegisterUser) (*response.RegisterUser, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.RegisterUser) *response.RegisterUser); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.RegisterUser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.RegisterUser) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: origCtx, payload, userId
func (_m *UsecaseCommand) UpdateUser(origCtx context.Context, payload request.UpdateUser, userId string) (string, error) {
	ret := _m.Called(origCtx, payload, userId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.UpdateUser, string) (string, error)); ok {
		return rf(origCtx, payload, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.UpdateUser, string) string); ok {
		r0 = rf(origCtx, payload, userId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.UpdateUser, string) error); ok {
		r1 = rf(origCtx, payload, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyRegisterUser provides a mock function with given fields: origCtx, payload
func (_m *UsecaseCommand) VerifyRegisterUser(origCtx context.Context, payload request.VerifyRegisterUser) (*response.VerifyRegister, error) {
	ret := _m.Called(origCtx, payload)

	if len(ret) == 0 {
		panic("no return value specified for VerifyRegisterUser")
	}

	var r0 *response.VerifyRegister
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.VerifyRegisterUser) (*response.VerifyRegister, error)); ok {
		return rf(origCtx, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.VerifyRegisterUser) *response.VerifyRegister); ok {
		r0 = rf(origCtx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.VerifyRegister)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.VerifyRegisterUser) error); ok {
		r1 = rf(origCtx, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsecaseCommand creates a new instance of UsecaseCommand. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecaseCommand(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsecaseCommand {
	mock := &UsecaseCommand{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
