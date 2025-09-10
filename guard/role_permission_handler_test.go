package guard

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/iam"
)

func TestRolePermissionHandler(t *testing.T) {
	tcs := map[string]struct {
		userProfile iam.UserProfile
		resource    string
		permissions Action
		mockEnforce bool
		enforceErr  error
		expErr      error
	}{
		"success": {
			userProfile: iam.NewUserProfile("user1", []string{"admin"}, nil),
			resource:    "resource1",
			permissions: ActionRead,
			mockEnforce: true,
			enforceErr:  nil,
			expErr:      nil,
		},
		"no user profile": {
			userProfile: iam.UserProfile{},
			resource:    "resource1",
			permissions: ActionRead,
			expErr:      errForbidden,
		},
		"enforcer action not allowed": {
			userProfile: iam.NewUserProfile("user1", []string{"user"}, nil),
			resource:    "resource1",
			permissions: ActionDelete,
			mockEnforce: true,
			enforceErr:  iam.ErrActionIsNotAllowed,
			expErr:      errForbidden,
		},
		"enforcer other error": {
			userProfile: iam.NewUserProfile("user1", []string{"user"}, nil),
			resource:    "resource1",
			permissions: ActionUpdate,
			mockEnforce: true,
			enforceErr:  errors.New("enforcer error"),
			expErr:      errors.New("enforcer error"),
		},
		"user with multiple roles": {
			userProfile: iam.NewUserProfile("user1", []string{"admin", "user"}, nil),
			resource:    "resource1",
			permissions: ActionCreate,
			mockEnforce: true,
			enforceErr:  nil,
			expErr:      nil,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			reqCtx := context.Background()
			request := httptest.NewRequestWithContext(reqCtx, http.MethodGet, "/", nil)
			request = request.WithContext(iam.SetUserProfileInContext(request.Context(), tc.userProfile))

			_, ctx, _ := lit.NewRouterForTest(httptest.NewRecorder())
			ctx.SetRequest(request)

			mockEnforcer := new(iam.MockEnforcer)
			if tc.mockEnforce {
				roles := tc.userProfile.GetRoles()
				var role string
				if len(roles) > 0 {
					role = roles[0]
				}
				mockEnforcer.On("Enforce", role, tc.resource, tc.permissions.String()).
					Return(tc.enforceErr)
			}

			guard := New(nil, mockEnforcer)

			handlerCalled := false
			mockHandler := func(c lit.Context) error {
				handlerCalled = true
				return nil
			}

			// When
			hdl := guard.RolePermissionHandler(mockHandler, tc.resource, tc.permissions)
			err := hdl(ctx)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.False(t, handlerCalled)
			} else {
				require.NoError(t, err)
				require.True(t, handlerCalled)
			}
			mockEnforcer.AssertExpectations(t)
		})
	}
}
