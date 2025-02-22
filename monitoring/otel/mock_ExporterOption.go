// Code generated by mockery v2.52.2. DO NOT EDIT.

package otel

import mock "github.com/stretchr/testify/mock"

// MockExporterOption is an autogenerated mock type for the ExporterOption type
type MockExporterOption struct {
	mock.Mock
}

type MockExporterOption_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExporterOption) EXPECT() *MockExporterOption_Expecter {
	return &MockExporterOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockExporterOption) Execute(_a0 *config) {
	_m.Called(_a0)
}

// MockExporterOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockExporterOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *config
func (_e *MockExporterOption_Expecter) Execute(_a0 interface{}) *MockExporterOption_Execute_Call {
	return &MockExporterOption_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockExporterOption_Execute_Call) Run(run func(_a0 *config)) *MockExporterOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*config))
	})
	return _c
}

func (_c *MockExporterOption_Execute_Call) Return() *MockExporterOption_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockExporterOption_Execute_Call) RunAndReturn(run func(*config)) *MockExporterOption_Execute_Call {
	_c.Run(run)
	return _c
}

// NewMockExporterOption creates a new instance of MockExporterOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExporterOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExporterOption {
	mock := &MockExporterOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
