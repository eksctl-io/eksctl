// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	mock "github.com/stretchr/testify/mock"
)

// ClusterStackDescriber is an autogenerated mock type for the ClusterStackDescriber type
type ClusterStackDescriber struct {
	mock.Mock
}

type ClusterStackDescriber_Expecter struct {
	mock *mock.Mock
}

func (_m *ClusterStackDescriber) EXPECT() *ClusterStackDescriber_Expecter {
	return &ClusterStackDescriber_Expecter{mock: &_m.Mock}
}

// ClusterHasDedicatedVPC provides a mock function with given fields: ctx
func (_m *ClusterStackDescriber) ClusterHasDedicatedVPC(ctx context.Context) (bool, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ClusterHasDedicatedVPC")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClusterStackDescriber_ClusterHasDedicatedVPC_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClusterHasDedicatedVPC'
type ClusterStackDescriber_ClusterHasDedicatedVPC_Call struct {
	*mock.Call
}

// ClusterHasDedicatedVPC is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ClusterStackDescriber_Expecter) ClusterHasDedicatedVPC(ctx interface{}) *ClusterStackDescriber_ClusterHasDedicatedVPC_Call {
	return &ClusterStackDescriber_ClusterHasDedicatedVPC_Call{Call: _e.mock.On("ClusterHasDedicatedVPC", ctx)}
}

func (_c *ClusterStackDescriber_ClusterHasDedicatedVPC_Call) Run(run func(ctx context.Context)) *ClusterStackDescriber_ClusterHasDedicatedVPC_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ClusterStackDescriber_ClusterHasDedicatedVPC_Call) Return(_a0 bool, _a1 error) *ClusterStackDescriber_ClusterHasDedicatedVPC_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClusterStackDescriber_ClusterHasDedicatedVPC_Call) RunAndReturn(run func(context.Context) (bool, error)) *ClusterStackDescriber_ClusterHasDedicatedVPC_Call {
	_c.Call.Return(run)
	return _c
}

// DescribeClusterStack provides a mock function with given fields: ctx
func (_m *ClusterStackDescriber) DescribeClusterStack(ctx context.Context) (*types.Stack, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DescribeClusterStack")
	}

	var r0 *types.Stack
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*types.Stack, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *types.Stack); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Stack)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClusterStackDescriber_DescribeClusterStack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DescribeClusterStack'
type ClusterStackDescriber_DescribeClusterStack_Call struct {
	*mock.Call
}

// DescribeClusterStack is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ClusterStackDescriber_Expecter) DescribeClusterStack(ctx interface{}) *ClusterStackDescriber_DescribeClusterStack_Call {
	return &ClusterStackDescriber_DescribeClusterStack_Call{Call: _e.mock.On("DescribeClusterStack", ctx)}
}

func (_c *ClusterStackDescriber_DescribeClusterStack_Call) Run(run func(ctx context.Context)) *ClusterStackDescriber_DescribeClusterStack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ClusterStackDescriber_DescribeClusterStack_Call) Return(_a0 *types.Stack, _a1 error) *ClusterStackDescriber_DescribeClusterStack_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClusterStackDescriber_DescribeClusterStack_Call) RunAndReturn(run func(context.Context) (*types.Stack, error)) *ClusterStackDescriber_DescribeClusterStack_Call {
	_c.Call.Return(run)
	return _c
}

// NewClusterStackDescriber creates a new instance of ClusterStackDescriber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClusterStackDescriber(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClusterStackDescriber {
	mock := &ClusterStackDescriber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
