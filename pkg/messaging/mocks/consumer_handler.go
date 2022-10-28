// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	mock "github.com/stretchr/testify/mock"
)

// ConsumerHandler is an autogenerated mock type for the ConsumerHandler type
type ConsumerHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: ctx, consumeContext
func (_m *ConsumerHandler) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	ret := _m.Called(ctx, consumeContext)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.MessageConsumeContext) error); ok {
		r0 = rf(ctx, consumeContext)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewConsumerHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewConsumerHandler creates a new instance of ConsumerHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConsumerHandler(t mockConstructorTestingTNewConsumerHandler) *ConsumerHandler {
	mock := &ConsumerHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
