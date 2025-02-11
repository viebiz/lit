package iam

import (
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	pkgerrors "github.com/pkg/errors"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/viebiz/lit/postgres"
)

const (
	permissionTableName = "permissions"
)

type permission struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100;uniqueIndex:unique_index"`
	V0    string `gorm:"size:100;uniqueIndex:unique_index"`
	V1    string `gorm:"size:100;uniqueIndex:unique_index"`
	V2    string `gorm:"size:100;uniqueIndex:unique_index"`
	V3    string `gorm:"size:100;uniqueIndex:unique_index"`
	V4    string `gorm:"size:100;uniqueIndex:unique_index"`
	V5    string `gorm:"size:100;uniqueIndex:unique_index"`
}

func newPostgresAdapter(db postgres.ContextExecutor) (persist.Adapter, error) {
	gormDB, err := gorm.Open(pgdriver.New(pgdriver.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	gormadapter.TurnOffAutoMigrate(gormDB)
	a, err := gormadapter.NewAdapterByDBWithCustomTable(gormDB, &permission{}, permissionTableName)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return a, nil
}
