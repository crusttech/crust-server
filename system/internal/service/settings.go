package service

import (
	"context"

	internalSettings "github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
)

type (
	// Wrapper service for system around internal settings service
	settings struct {
		db  db
		ctx context.Context

		internalSettings internalSettings.Service
	}

	SettingsService interface {
		With(ctx context.Context) SettingsService
		FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error)
		Set(v *internalSettings.Value) (err error)
		Get(name string, ownedBy uint64) (out *internalSettings.Value, err error)
	}
)

func Settings() SettingsService {
	return (&settings{}).With(context.Background())
}

func (svc settings) With(ctx context.Context) SettingsService {
	db := repository.DB(ctx)
	return &settings{
		ctx:              ctx,
		internalSettings: internalSettings.NewService(internalSettings.NewRepository(db, "sys_settings")).With(ctx),
	}
}

func (svc settings) FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error) {
	return svc.internalSettings.FindByPrefix(prefix)
}

func (svc settings) Set(v *internalSettings.Value) (err error) {
	return svc.internalSettings.Set(v)
}

func (svc settings) Get(name string, ownedBy uint64) (out *internalSettings.Value, err error) {
	return svc.internalSettings.Get(name, ownedBy)
}
