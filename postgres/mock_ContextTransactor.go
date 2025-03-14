// Code generated by mockery v2.52.2. DO NOT EDIT.

package postgres

import (
	context "context"
	sql "database/sql"

	mock "github.com/stretchr/testify/mock"
)

// MockContextTransactor is an autogenerated mock type for the ContextTransactor type
type MockContextTransactor struct {
	mock.Mock
}

type MockContextTransactor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockContextTransactor) EXPECT() *MockContextTransactor_Expecter {
	return &MockContextTransactor_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with no fields
func (_m *MockContextTransactor) Commit() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockContextTransactor_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type MockContextTransactor_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
func (_e *MockContextTransactor_Expecter) Commit() *MockContextTransactor_Commit_Call {
	return &MockContextTransactor_Commit_Call{Call: _e.mock.On("Commit")}
}

func (_c *MockContextTransactor_Commit_Call) Run(run func()) *MockContextTransactor_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockContextTransactor_Commit_Call) Return(_a0 error) *MockContextTransactor_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContextTransactor_Commit_Call) RunAndReturn(run func() error) *MockContextTransactor_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Exec provides a mock function with given fields: query, args
func (_m *MockContextTransactor) Exec(query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (sql.Result, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) sql.Result); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockContextTransactor_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) Exec(query interface{}, args ...interface{}) *MockContextTransactor_Exec_Call {
	return &MockContextTransactor_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{query}, args...)...)}
}

func (_c *MockContextTransactor_Exec_Call) Run(run func(query string, args ...interface{})) *MockContextTransactor_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_Exec_Call) Return(_a0 sql.Result, _a1 error) *MockContextTransactor_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_Exec_Call) RunAndReturn(run func(string, ...interface{}) (sql.Result, error)) *MockContextTransactor_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// ExecContext provides a mock function with given fields: ctx, query, args
func (_m *MockContextTransactor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ExecContext")
	}

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (sql.Result, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) sql.Result); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_ExecContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecContext'
type MockContextTransactor_ExecContext_Call struct {
	*mock.Call
}

// ExecContext is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) ExecContext(ctx interface{}, query interface{}, args ...interface{}) *MockContextTransactor_ExecContext_Call {
	return &MockContextTransactor_ExecContext_Call{Call: _e.mock.On("ExecContext",
		append([]interface{}{ctx, query}, args...)...)}
}

func (_c *MockContextTransactor_ExecContext_Call) Run(run func(ctx context.Context, query string, args ...interface{})) *MockContextTransactor_ExecContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_ExecContext_Call) Return(_a0 sql.Result, _a1 error) *MockContextTransactor_ExecContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_ExecContext_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (sql.Result, error)) *MockContextTransactor_ExecContext_Call {
	_c.Call.Return(run)
	return _c
}

// Prepare provides a mock function with given fields: query
func (_m *MockContextTransactor) Prepare(query string) (*sql.Stmt, error) {
	ret := _m.Called(query)

	if len(ret) == 0 {
		panic("no return value specified for Prepare")
	}

	var r0 *sql.Stmt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*sql.Stmt, error)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(string) *sql.Stmt); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_Prepare_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Prepare'
type MockContextTransactor_Prepare_Call struct {
	*mock.Call
}

// Prepare is a helper method to define mock.On call
//   - query string
func (_e *MockContextTransactor_Expecter) Prepare(query interface{}) *MockContextTransactor_Prepare_Call {
	return &MockContextTransactor_Prepare_Call{Call: _e.mock.On("Prepare", query)}
}

func (_c *MockContextTransactor_Prepare_Call) Run(run func(query string)) *MockContextTransactor_Prepare_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockContextTransactor_Prepare_Call) Return(_a0 *sql.Stmt, _a1 error) *MockContextTransactor_Prepare_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_Prepare_Call) RunAndReturn(run func(string) (*sql.Stmt, error)) *MockContextTransactor_Prepare_Call {
	_c.Call.Return(run)
	return _c
}

// PrepareContext provides a mock function with given fields: ctx, query
func (_m *MockContextTransactor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for PrepareContext")
	}

	var r0 *sql.Stmt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*sql.Stmt, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *sql.Stmt); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_PrepareContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrepareContext'
type MockContextTransactor_PrepareContext_Call struct {
	*mock.Call
}

// PrepareContext is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
func (_e *MockContextTransactor_Expecter) PrepareContext(ctx interface{}, query interface{}) *MockContextTransactor_PrepareContext_Call {
	return &MockContextTransactor_PrepareContext_Call{Call: _e.mock.On("PrepareContext", ctx, query)}
}

func (_c *MockContextTransactor_PrepareContext_Call) Run(run func(ctx context.Context, query string)) *MockContextTransactor_PrepareContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockContextTransactor_PrepareContext_Call) Return(_a0 *sql.Stmt, _a1 error) *MockContextTransactor_PrepareContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_PrepareContext_Call) RunAndReturn(run func(context.Context, string) (*sql.Stmt, error)) *MockContextTransactor_PrepareContext_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: query, args
func (_m *MockContextTransactor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockContextTransactor_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) Query(query interface{}, args ...interface{}) *MockContextTransactor_Query_Call {
	return &MockContextTransactor_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{query}, args...)...)}
}

func (_c *MockContextTransactor_Query_Call) Run(run func(query string, args ...interface{})) *MockContextTransactor_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_Query_Call) Return(_a0 *sql.Rows, _a1 error) *MockContextTransactor_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_Query_Call) RunAndReturn(run func(string, ...interface{}) (*sql.Rows, error)) *MockContextTransactor_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryContext provides a mock function with given fields: ctx, query, args
func (_m *MockContextTransactor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryContext")
	}

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Rows); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContextTransactor_QueryContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryContext'
type MockContextTransactor_QueryContext_Call struct {
	*mock.Call
}

// QueryContext is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) QueryContext(ctx interface{}, query interface{}, args ...interface{}) *MockContextTransactor_QueryContext_Call {
	return &MockContextTransactor_QueryContext_Call{Call: _e.mock.On("QueryContext",
		append([]interface{}{ctx, query}, args...)...)}
}

func (_c *MockContextTransactor_QueryContext_Call) Run(run func(ctx context.Context, query string, args ...interface{})) *MockContextTransactor_QueryContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_QueryContext_Call) Return(_a0 *sql.Rows, _a1 error) *MockContextTransactor_QueryContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContextTransactor_QueryContext_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (*sql.Rows, error)) *MockContextTransactor_QueryContext_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: query, args
func (_m *MockContextTransactor) QueryRow(query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRow")
	}

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Row); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// MockContextTransactor_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type MockContextTransactor_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) QueryRow(query interface{}, args ...interface{}) *MockContextTransactor_QueryRow_Call {
	return &MockContextTransactor_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{query}, args...)...)}
}

func (_c *MockContextTransactor_QueryRow_Call) Run(run func(query string, args ...interface{})) *MockContextTransactor_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_QueryRow_Call) Return(_a0 *sql.Row) *MockContextTransactor_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContextTransactor_QueryRow_Call) RunAndReturn(run func(string, ...interface{}) *sql.Row) *MockContextTransactor_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRowContext provides a mock function with given fields: ctx, query, args
func (_m *MockContextTransactor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRowContext")
	}

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Row); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// MockContextTransactor_QueryRowContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRowContext'
type MockContextTransactor_QueryRowContext_Call struct {
	*mock.Call
}

// QueryRowContext is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
//   - args ...interface{}
func (_e *MockContextTransactor_Expecter) QueryRowContext(ctx interface{}, query interface{}, args ...interface{}) *MockContextTransactor_QueryRowContext_Call {
	return &MockContextTransactor_QueryRowContext_Call{Call: _e.mock.On("QueryRowContext",
		append([]interface{}{ctx, query}, args...)...)}
}

func (_c *MockContextTransactor_QueryRowContext_Call) Run(run func(ctx context.Context, query string, args ...interface{})) *MockContextTransactor_QueryRowContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockContextTransactor_QueryRowContext_Call) Return(_a0 *sql.Row) *MockContextTransactor_QueryRowContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContextTransactor_QueryRowContext_Call) RunAndReturn(run func(context.Context, string, ...interface{}) *sql.Row) *MockContextTransactor_QueryRowContext_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with no fields
func (_m *MockContextTransactor) Rollback() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockContextTransactor_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type MockContextTransactor_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
func (_e *MockContextTransactor_Expecter) Rollback() *MockContextTransactor_Rollback_Call {
	return &MockContextTransactor_Rollback_Call{Call: _e.mock.On("Rollback")}
}

func (_c *MockContextTransactor_Rollback_Call) Run(run func()) *MockContextTransactor_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockContextTransactor_Rollback_Call) Return(_a0 error) *MockContextTransactor_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContextTransactor_Rollback_Call) RunAndReturn(run func() error) *MockContextTransactor_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockContextTransactor creates a new instance of MockContextTransactor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockContextTransactor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockContextTransactor {
	mock := &MockContextTransactor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
