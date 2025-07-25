// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// NodeGroupResourceSet is an autogenerated mock type for the NodeGroupResourceSet type
type NodeGroupResourceSet struct {
	mock.Mock
}

// AddAllResources provides a mock function with given fields: ctx
func (_m *NodeGroupResourceSet) AddAllResources(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for AddAllResources")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllOutputs provides a mock function with given fields: _a0
func (_m *NodeGroupResourceSet) GetAllOutputs(_a0 types.Stack) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetAllOutputs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Stack) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenderJSON provides a mock function with no fields
func (_m *NodeGroupResourceSet) RenderJSON() ([]byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RenderJSON")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithIAM provides a mock function with no fields
func (_m *NodeGroupResourceSet) WithIAM() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WithIAM")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// WithNamedIAM provides a mock function with no fields
func (_m *NodeGroupResourceSet) WithNamedIAM() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WithNamedIAM")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewNodeGroupResourceSet creates a new instance of NodeGroupResourceSet. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNodeGroupResourceSet(t interface {
	mock.TestingT
	Cleanup(func())
}) *NodeGroupResourceSet {
	mock := &NodeGroupResourceSet{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
