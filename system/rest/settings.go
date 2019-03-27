package rest

import (
	"context"
	"github.com/crusttech/crust/system/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Settings struct{}

func (Settings) New() *Settings {
	return &Settings{}
}

func (ctrl *Settings) List(ctx context.Context, r *request.SettingsList) (interface{}, error) {
	return nil, errors.New("Not implemented: Settings.list")
}

func (ctrl *Settings) Get(ctx context.Context, r *request.SettingsGet) (interface{}, error) {
	return nil, errors.New("Not implemented: Settings.get")
}

func (ctrl *Settings) Set(ctx context.Context, r *request.SettingsSet) (interface{}, error) {
	return nil, errors.New("Not implemented: Settings.set")
}
