package service

import (
	"github.com/crusttech/crust/system/internal/service"
)

func Init() error {
	err := service.Init()
	DefaultRules = service.DefaultRules
	DefaultUser = service.DefaultUser
	return err
}
