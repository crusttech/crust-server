package rest

import (
	"context"

	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Settings struct {
		svc struct {
			settings service.SettingsService
		}
	}
)

func (Settings) New() *Settings {
	ctrl := &Settings{}
	ctrl.svc.settings = service.DefaultSettings

	return ctrl
}

func (ctrl *Settings) List(ctx context.Context, r *request.SettingsList) (interface{}, error) {
	if vv, err := ctrl.svc.settings.FindByPrefix(r.Prefix); err != nil {
		return nil, err
	} else {
		return vv, err
	}
}

func (ctrl *Settings) Update(ctx context.Context, r *request.SettingsUpdate) (interface{}, error) {
	return nil, errors.New("Not implemented: Settings.update")
}

func (ctrl *Settings) Get(ctx context.Context, r *request.SettingsGet) (interface{}, error) {
	if v, err := ctrl.svc.settings.With(ctx).Get(r.Key, r.OwnerID); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

func (ctrl *Settings) Set(ctx context.Context, r *request.SettingsSet) (interface{}, error) {
	return nil, errors.New("Not implemented: Settings.set")
}
