package testutil

import (
	"github.com/google/go-cmp/cmp"
)

// Equal compares the expected and actual values of type T.
//
// It fails the test immediately if the values are not equal,
// and reports the difference using cmp.Diff.
//
// Options can be provided via opts to customize the comparison.
func Equal[T any](t TestingT, expected, actual T, opts ...Option[T]) {
	t.Helper()

	cmpOpts := make([]cmp.Option, len(opts))
	for i, opt := range opts {
		opt.check(expected)
		cmpOpts[i] = opt.toCmpOption()
	}

	if !cmp.Equal(expected, actual, cmpOpts...) {
		t.Errorf("\n mismatched. \n expected: %+v \n got: %+v \n diff:\n%s",
			expected, actual,
			cmp.Diff(expected, actual, cmpOpts...))
		t.FailNow()
	}
}
