package service

import (
	"context"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/auth"
	internalRules "github.com/crusttech/crust/internal/rules"
	systemService "github.com/crusttech/crust/system/service"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		rules systemService.RulesService
	}

	resource interface {
		Resource() internalRules.Resource
	}

	PermissionsService interface {
		With(context.Context) PermissionsService

		CanAccessCompose() bool
	}

	Compose struct{}
)

func (Compose) Resource() internalRules.Resource {
	return internalRules.Resource{Service: "compose"}
}

func Permissions() PermissionsService {
	return (&permissions{
		rules: systemService.DefaultRules,
	}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		rules: p.rules.With(ctx),
	}
}

func (p *permissions) baseResource() resource {
	return &Compose{}
}

func (p *permissions) CanAccessCompose() bool {
	return p.checkAccess(p.baseResource(), "access")
}

func (p *permissions) CanCreateNamspace() bool {
	return p.checkAccess(p.baseResource(), "namespace.create")
}

func (p *permissions) CanCreateModule() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: "crm"}
	return p.checkAccess(ns, "module.create")
}

func (p *permissions) CanReadModule(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateModule(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteModule(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreateRecord(r resource) bool {
	return p.checkAccess(r, "record.create")
}

func (p *permissions) CanReadRecord(r *types.Record) bool {
	return p.checkAccess(r, "record.read", p.recordOwnerFallback(r))
}

func (p *permissions) CanUpdateRecord(r *types.Record) bool {
	return p.checkAccess(r, "record.update", p.recordOwnerFallback(r))
}

func (p *permissions) canDeleteRecord(r *types.Record) bool {
	return p.checkAccess(r, "record.delete", p.recordOwnerFallback(r))
}

func (p permissions) recordOwnerFallback(r *types.Record) func() internalRules.Access {
	return func() internalRules.Access {
		if auth.GetIdentityFromContext(p.ctx).Identity() == r.OwnedBy {
			return internalRules.Allow
		}

		return internalRules.Deny
	}
}

func (p *permissions) CanCreateChart() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: "crm"}
	return p.checkAccess(ns, "chart.create")
}

func (p *permissions) CanReadChart(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateChart(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteChart(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreateTrigger() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: "crm"}
	return p.checkAccess(ns, "trigger.create")
}

func (p *permissions) CanReadTrigger(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdateTrigger(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeleteTrigger(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) CanCreatePage() bool {
	// @todo move to func args when namespaces are implemented
	ns := &types.Namespace{ID: "crm"}
	return p.checkAccess(ns, "page.create")
}

func (p *permissions) CanReadPage(r resource) bool {
	return p.checkAccess(r, "read")
}

func (p *permissions) CanUpdatePage(r resource) bool {
	return p.checkAccess(r, "update")
}

func (p *permissions) CanDeletePage(r resource) bool {
	return p.checkAccess(r, "delete")
}

func (p *permissions) checkAccess(resource resource, operation string, fallbacks ...internalRules.CheckAccessFunc) bool {
	access := p.rules.Check(resource.Resource().String(), operation, fallbacks...)
	if access == internalRules.Allow {
		return true
	}
	return false
}
