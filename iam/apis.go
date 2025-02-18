package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"

	pkgerrors "github.com/pkg/errors"

	"github.com/viebiz/lit/jwt"
	"github.com/viebiz/lit/monitoring"
)

func NewEnforcer(ctx context.Context, cfg EnforcerConfig) (Enforcer, error) {
	// 1. Read and parse model config
	m, err := model.NewModelFromString(authModel)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	adapter, err := newPostgresAdapter(cfg.DBConn)
	if err != nil {
		return nil, err
	}

	logger := enforcerLogger{
		Logger: monitoring.FromContext(ctx),
	}

	// 2. Init casbin enforcer with model
	cbEnforcer, err := casbin.NewEnforcer(m, adapter, &logger)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	cbEnforcer.AddFunction(hasPermissionKeyMatch, hasPermission)

	return enforcer{
		cb: cbEnforcer,
	}, nil
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

func NewRFC9068Validator(issuer, audience string, client HTTPClient) (Validator, error) {
	jwksURI := fmt.Sprintf("%s/.well-known/jwks.json", strings.TrimSuffix(issuer, "/"))
	v := rfc9068Validator{
		jwksURI:     jwksURI,
		issuer:      issuer,
		audience:    audience,
		httpClient:  client,
		tokenParser: jwt.NewParser[Claims](),
	}

	// Download signing key
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := v.downloadSigningKey(ctx); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewUserProfile(id string, roles []string, permissions []string) UserProfile {
	return UserProfile{
		id:          id,
		roles:       roles,
		permissions: permissions,
	}
}

func ExtractUserProfileFromClaims(claims Claims) (UserProfile, error) {
	sub := claims.RegisteredClaims.Subject

	roles, err := extractRolesFromClaims(claims)
	if err != nil {
		return UserProfile{}, err
	}

	return UserProfile{
		id:    sub,
		roles: roles,
	}, nil
}
