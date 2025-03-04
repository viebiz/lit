package testutil

type TestingT interface {
	Helper()

	Errorf(format string, args ...any)

	FailNow()
}
