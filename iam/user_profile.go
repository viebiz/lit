package iam

import (
	"fmt"
	"slices"
	"strings"
)

const (
	roleClaimKey       string = "roles"
	permissionClaimKey string = "permissions"
)

type UserProfile struct {
	id          string
	roles       []string
	permissions []string
}

func (p UserProfile) ID() string {
	return p.id
}

func (p UserProfile) GetRoles() []string {
	return slices.Clone(p.roles)
}

func (p UserProfile) GetPermission() []string {
	return slices.Clone(p.permissions)
}

func (p UserProfile) GetRoleString() string {
	return strings.Join(p.roles, ",")
}

func extractRolesFromClaims(claims Claims) ([]string, error) {
	rolesClaim, exists := claims.ExtraClaims[roleClaimKey]
	if !exists {
		return nil, ErrMissingRequiredClaim
	}

	switch v := rolesClaim.(type) {
	case string:
		return strings.Split(v, ","), nil
	case []string:
		return v, nil
	case []interface{}:
		rs := make([]string, len(v))
		for idx, item := range v {
			role, ok := item.(string)
			if !ok {
				role = fmt.Sprintf("%s", item)
			}

			rs[idx] = role
		}

		return rs, nil
	default:
		return nil, ErrInvalidToken
	}
}
