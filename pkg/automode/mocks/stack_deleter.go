// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	mock "github.com/stretchr/testify/mock"
)

// StackDeleter is an autogenerated mock type for the StackDeleter type
type StackDeleter struct {
	mock.Mock
}

type StackDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *StackDeleter) EXPECT() *StackDeleter_Expecter {
	return &StackDeleter_Expecter{mock: &_m.Mock}
}

// DeleteStackSync provides a mock function with given fields: _a0, _a1
func (_m *StackDeleter) DeleteStackSync(_a0 context.Context, _a1 *types.Stack) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteStackSync")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StackDeleter_DeleteStackSync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteStackSync'
type StackDeleter_DeleteStackSync_Call struct {
	*mock.Call
}

// DeleteStackSync is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *types.Stack
func (_e *StackDeleter_Expecter) DeleteStackSync(_a0 interface{}, _a1 interface{}) *StackDeleter_DeleteStackSync_Call {
	return &StackDeleter_DeleteStackSync_Call{Call: _e.mock.On("DeleteStackSync", _a0, _a1)}
}

func (_c *StackDeleter_DeleteStackSync_Call) Run(run func(_a0 context.Context, _a1 *types.Stack)) *StackDeleter_DeleteStackSync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.Stack))
	})
	return _c
}

func (_c *StackDeleter_DeleteStackSync_Call) Return(_a0 error) *StackDeleter_DeleteStackSync_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StackDeleter_DeleteStackSync_Call) RunAndReturn(run func(context.Context, *types.Stack) error) *StackDeleter_DeleteStackSync_Call {
	_c.Call.Return(run)
	return _c
}

// DescribeStack provides a mock function with given fields: ctx, stack
func (_m *StackDeleter) DescribeStack(ctx context.Context, stack *types.Stack) (*types.Stack, error) {
	ret := _m.Called(ctx, stack)

	if len(ret) == 0 {
		panic("no return value specified for DescribeStack")
	}

	var r0 *types.Stack
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack) (*types.Stack, error)); ok {
		return rf(ctx, stack)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.Stack) *types.Stack); ok {
		r0 = rf(ctx, stack)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Stack)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.Stack) error); ok {
		r1 = rf(ctx, stack)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StackDeleter_DescribeStack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DescribeStack'
type StackDeleter_DescribeStack_Call struct {
	*mock.Call
}

// DescribeStack is a helper method to define mock.On call
//   - ctx context.Context
//   - stack *types.Stack
func (_e *StackDeleter_Expecter) DescribeStack(ctx interface{}, stack interface{}) *StackDeleter_DescribeStack_Call {
	return &StackDeleter_DescribeStack_Call{Call: _e.mock.On("DescribeStack", ctx, stack)}
}

func (_c *StackDeleter_DescribeStack_Call) Run(run func(ctx context.Context, stack *types.Stack)) *StackDeleter_DescribeStack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.Stack))
	})
	return _c
}

func (_c *StackDeleter_DescribeStack_Call) Return(_a0 *types.Stack, _a1 error) *StackDeleter_DescribeStack_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackDeleter_DescribeStack_Call) RunAndReturn(run func(context.Context, *types.Stack) (*types.Stack, error)) *StackDeleter_DescribeStack_Call {
	_c.Call.Return(run)
	return _c
}

// NewStackDeleter creates a new instance of StackDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStackDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *StackDeleter {
	mock := &StackDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
