package guard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAction_String(t *testing.T) {
	tcs := map[string]struct {
		action Action
		exp    string
	}{
		"read": {
			action: ActionRead,
			exp:    "R",
		},
		"create": {
			action: ActionCreate,
			exp:    "C",
		},
		"update": {
			action: ActionUpdate,
			exp:    "U",
		},
		"delete": {
			action: ActionDelete,
			exp:    "D",
		},
		"custom": {
			action: Action("X"),
			exp:    "X",
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			result := tc.action.String()
			require.Equal(t, tc.exp, result)
		})
	}
}

func TestAction_IsValid(t *testing.T) {
	tcs := map[string]struct {
		action Action
		exp    bool
	}{
		"read": {
			action: ActionRead,
			exp:    true,
		},
		"create": {
			action: ActionCreate,
			exp:    true,
		},
		"update": {
			action: ActionUpdate,
			exp:    true,
		},
		"delete": {
			action: ActionDelete,
			exp:    true,
		},
		"invalid": {
			action: Action("X"),
			exp:    false,
		},
		"empty": {
			action: Action(""),
			exp:    false,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			result := tc.action.IsValid()
			require.Equal(t, tc.exp, result)
		})
	}
}
