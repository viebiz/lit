package iam

import (
	"fmt"
	"strings"
)

const (
	scopeClaimKey  string = "scope"
	scopeSeparator string = " "
)

type M2MProfile struct {
	id     string
	scopes map[string]bool
}

func NewM2MProfile(id string, scopes []string) M2MProfile {
	scopeMap := make(map[string]bool)
	for _, scope := range scopes {
		scopeMap[scope] = true
	}

	return M2MProfile{
		id:     id,
		scopes: scopeMap,
	}
}

func ExtractM2MProfileFromClaims(claims Claims) (M2MProfile, error) {
	sub := claims.RegisteredClaims.Subject

	scopeSet, err := extractScopeFromClaims(claims)
	if err != nil {
		return M2MProfile{}, err
	}

	return M2MProfile{
		id:     sub,
		scopes: scopeSet,
	}, nil
}

func (p M2MProfile) ID() string {
	return p.id
}

func (p M2MProfile) GetScopes() []string {
	scopes := make([]string, 0, len(p.scopes))
	for scope := range p.scopes {
		scopes = append(scopes, scope)
	}

	return scopes
}

func (p M2MProfile) HasScope(scope string) bool {
	if match, exists := p.scopes[scope]; exists {
		return match
	}

	return false
}

func (p M2MProfile) HasAnyScope(scopes ...string) bool {
	for _, s := range scopes {
		if p.HasScope(s) {
			return true
		}
	}

	return false
}

func extractScopeFromClaims(claims Claims) (map[string]bool, error) {
	scopeClaim, exists := claims.ExtraClaims[scopeClaimKey]
	if !exists {
		return nil, ErrMissingRequiredClaim
	}

	scopes, ok := scopeClaim.(string)
	if !ok {
		scopes = fmt.Sprintf("%s", scopeClaim)
	}

	scopeSet := make(map[string]bool)
	for _, scope := range strings.Split(scopes, scopeSeparator) {
		scopeSet[scope] = true
	}

	return scopeSet, nil
}
