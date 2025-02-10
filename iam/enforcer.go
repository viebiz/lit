package iam

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	pkgerrors "github.com/pkg/errors"

	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/postgres"
)

type Enforcer interface {
	Enforce(sub, obj, act string) error
}

type enforcer struct {
	cb casbin.IEnforcer
}

type EnforcerConfig struct {
	DBConn postgres.ContextExecutor
}

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
