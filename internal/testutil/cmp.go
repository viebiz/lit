package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal[T any](t *testing.T, expected, actual T, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(expected, actual, opts...) {
		t.Errorf("\n mismatched. \n expected: %+v \n got: %+v \n diff:\n%s",
			expected, actual,
			cmp.Diff(expected, actual, opts...))
		t.FailNow()
	}
}
