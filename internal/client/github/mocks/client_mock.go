// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	github "bot/internal/client/github"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ClientMock is an autogenerated mock type for the ClientMock type
type ClientMock struct {
	mock.Mock
}

type ClientMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ClientMock) EXPECT() *ClientMock_Expecter {
	return &ClientMock_Expecter{mock: &_m.Mock}
}

// ChatCompletions provides a mock function with given fields: ctx, request
func (_m *ClientMock) ChatCompletions(ctx context.Context, request *github.ChatCompletionRequest) (*github.ChatCompletionResponse, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for ChatCompletions")
	}

	var r0 *github.ChatCompletionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *github.ChatCompletionRequest) (*github.ChatCompletionResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *github.ChatCompletionRequest) *github.ChatCompletionResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.ChatCompletionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *github.ChatCompletionRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientMock_ChatCompletions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChatCompletions'
type ClientMock_ChatCompletions_Call struct {
	*mock.Call
}

// ChatCompletions is a helper method to define mock.On call
//   - ctx context.Context
//   - request *github.ChatCompletionRequest
func (_e *ClientMock_Expecter) ChatCompletions(ctx interface{}, request interface{}) *ClientMock_ChatCompletions_Call {
	return &ClientMock_ChatCompletions_Call{Call: _e.mock.On("ChatCompletions", ctx, request)}
}

func (_c *ClientMock_ChatCompletions_Call) Run(run func(ctx context.Context, request *github.ChatCompletionRequest)) *ClientMock_ChatCompletions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*github.ChatCompletionRequest))
	})
	return _c
}

func (_c *ClientMock_ChatCompletions_Call) Return(_a0 *github.ChatCompletionResponse, _a1 error) *ClientMock_ChatCompletions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClientMock_ChatCompletions_Call) RunAndReturn(run func(context.Context, *github.ChatCompletionRequest) (*github.ChatCompletionResponse, error)) *ClientMock_ChatCompletions_Call {
	_c.Call.Return(run)
	return _c
}

// NewClientMock creates a new instance of ClientMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientMock {
	mock := &ClientMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
