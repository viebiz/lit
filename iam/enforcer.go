package iam

import (
	"github.com/casbin/casbin/v2"
	pkgerrors "github.com/pkg/errors"

	"github.com/viebiz/lit/postgres"
)

type enforcer struct {
	cb casbin.IEnforcer
}

type EnforcerConfig struct {
	DBConn postgres.ContextExecutor
}

func (e enforcer) Enforce(sub, obj, act string) error {
	allowed, err := e.cb.Enforce(sub, obj, act)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	if !allowed {
		return ErrActionIsNotAllowed
	}

	return nil
}
