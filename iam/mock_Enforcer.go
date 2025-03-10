// Code generated by mockery v2.52.2. DO NOT EDIT.

package iam

import mock "github.com/stretchr/testify/mock"

// MockEnforcer is an autogenerated mock type for the Enforcer type
type MockEnforcer struct {
	mock.Mock
}

type MockEnforcer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEnforcer) EXPECT() *MockEnforcer_Expecter {
	return &MockEnforcer_Expecter{mock: &_m.Mock}
}

// Enforce provides a mock function with given fields: sub, obj, act
func (_m *MockEnforcer) Enforce(sub string, obj string, act string) error {
	ret := _m.Called(sub, obj, act)

	if len(ret) == 0 {
		panic("no return value specified for Enforce")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(sub, obj, act)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEnforcer_Enforce_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Enforce'
type MockEnforcer_Enforce_Call struct {
	*mock.Call
}

// Enforce is a helper method to define mock.On call
//   - sub string
//   - obj string
//   - act string
func (_e *MockEnforcer_Expecter) Enforce(sub interface{}, obj interface{}, act interface{}) *MockEnforcer_Enforce_Call {
	return &MockEnforcer_Enforce_Call{Call: _e.mock.On("Enforce", sub, obj, act)}
}

func (_c *MockEnforcer_Enforce_Call) Run(run func(sub string, obj string, act string)) *MockEnforcer_Enforce_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockEnforcer_Enforce_Call) Return(_a0 error) *MockEnforcer_Enforce_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEnforcer_Enforce_Call) RunAndReturn(run func(string, string, string) error) *MockEnforcer_Enforce_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEnforcer creates a new instance of MockEnforcer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEnforcer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEnforcer {
	mock := &MockEnforcer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
