// Code generated by mockery v2.52.2. DO NOT EDIT.

package jwt

import mock "github.com/stretchr/testify/mock"

// MockVerifyKey is an autogenerated mock type for the VerifyKey type
type MockVerifyKey struct {
	mock.Mock
}

type MockVerifyKey_Expecter struct {
	mock *mock.Mock
}

func (_m *MockVerifyKey) EXPECT() *MockVerifyKey_Expecter {
	return &MockVerifyKey_Expecter{mock: &_m.Mock}
}

// NewMockVerifyKey creates a new instance of MockVerifyKey. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockVerifyKey(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockVerifyKey {
	mock := &MockVerifyKey{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
