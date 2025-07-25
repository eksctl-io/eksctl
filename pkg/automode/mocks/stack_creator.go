// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	builder "github.com/weaveworks/eksctl/pkg/cfn/builder"

	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// StackCreator is an autogenerated mock type for the StackCreator type
type StackCreator struct {
	mock.Mock
}

type StackCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *StackCreator) EXPECT() *StackCreator_Expecter {
	return &StackCreator_Expecter{mock: &_m.Mock}
}

// CreateStack provides a mock function with given fields: ctx, stackName, resourceSet, tags, parameters, errs
func (_m *StackCreator) CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags map[string]string, parameters map[string]string, errs chan error) error {
	ret := _m.Called(ctx, stackName, resourceSet, tags, parameters, errs)

	if len(ret) == 0 {
		panic("no return value specified for CreateStack")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, builder.ResourceSetReader, map[string]string, map[string]string, chan error) error); ok {
		r0 = rf(ctx, stackName, resourceSet, tags, parameters, errs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StackCreator_CreateStack_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateStack'
type StackCreator_CreateStack_Call struct {
	*mock.Call
}

// CreateStack is a helper method to define mock.On call
//   - ctx context.Context
//   - stackName string
//   - resourceSet builder.ResourceSetReader
//   - tags map[string]string
//   - parameters map[string]string
//   - errs chan error
func (_e *StackCreator_Expecter) CreateStack(ctx interface{}, stackName interface{}, resourceSet interface{}, tags interface{}, parameters interface{}, errs interface{}) *StackCreator_CreateStack_Call {
	return &StackCreator_CreateStack_Call{Call: _e.mock.On("CreateStack", ctx, stackName, resourceSet, tags, parameters, errs)}
}

func (_c *StackCreator_CreateStack_Call) Run(run func(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags map[string]string, parameters map[string]string, errs chan error)) *StackCreator_CreateStack_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(builder.ResourceSetReader), args[3].(map[string]string), args[4].(map[string]string), args[5].(chan error))
	})
	return _c
}

func (_c *StackCreator_CreateStack_Call) Return(_a0 error) *StackCreator_CreateStack_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StackCreator_CreateStack_Call) RunAndReturn(run func(context.Context, string, builder.ResourceSetReader, map[string]string, map[string]string, chan error) error) *StackCreator_CreateStack_Call {
	_c.Call.Return(run)
	return _c
}

// GetClusterStackIfExists provides a mock function with given fields: ctx
func (_m *StackCreator) GetClusterStackIfExists(ctx context.Context) (*types.Stack, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetClusterStackIfExists")
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

// StackCreator_GetClusterStackIfExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetClusterStackIfExists'
type StackCreator_GetClusterStackIfExists_Call struct {
	*mock.Call
}

// GetClusterStackIfExists is a helper method to define mock.On call
//   - ctx context.Context
func (_e *StackCreator_Expecter) GetClusterStackIfExists(ctx interface{}) *StackCreator_GetClusterStackIfExists_Call {
	return &StackCreator_GetClusterStackIfExists_Call{Call: _e.mock.On("GetClusterStackIfExists", ctx)}
}

func (_c *StackCreator_GetClusterStackIfExists_Call) Run(run func(ctx context.Context)) *StackCreator_GetClusterStackIfExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *StackCreator_GetClusterStackIfExists_Call) Return(_a0 *types.Stack, _a1 error) *StackCreator_GetClusterStackIfExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StackCreator_GetClusterStackIfExists_Call) RunAndReturn(run func(context.Context) (*types.Stack, error)) *StackCreator_GetClusterStackIfExists_Call {
	_c.Call.Return(run)
	return _c
}

// NewStackCreator creates a new instance of StackCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStackCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *StackCreator {
	mock := &StackCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
