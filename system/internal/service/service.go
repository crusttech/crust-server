package service

import (
	"sync"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                   sync.Once
	DefaultSettings     SettingsService
	DefaultAuth         AuthService
	DefaultUser         UserService
	DefaultRole         RoleService
	DefaultRules        RulesService
	DefaultOrganisation OrganisationService
	DefaultApplication  ApplicationService
	DefaultPermissions  PermissionsService
)

func Init() {
	o.Do(func() {
		DefaultRules = Rules()
		DefaultSettings = Settings()
		DefaultPermissions = Permissions()
		DefaultAuth = Auth()
		DefaultUser = User()
		DefaultRole = Role()
		DefaultOrganisation = Organisation()
		DefaultApplication = Application()
	})
}
