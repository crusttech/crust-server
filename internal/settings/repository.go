package settings

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"
)

type (
	repository struct {
		dbh *factory.DB

		// sql table reference
		dbTable string
	}

	Repository interface {
		With(ctx context.Context) Repository

		Find(filter Filter) (ss ValueSet, err error)

		Get(name string, ownedBy uint64) (value *Value, err error)
		Set(value *Value) error
	}
)

func NewRepository(db *factory.DB, table string) Repository {
	return &repository{
		dbTable: table,
		dbh:     db,
	}
}

func (r *repository) db() *factory.DB {
	return r.dbh
}

func (r *repository) With(ctx context.Context) Repository {
	return &repository{
		dbTable: r.dbTable,
		dbh:     r.db().With(ctx),
	}
}

func (r *repository) Find(f Filter) (ss ValueSet, err error) {
	f.Normalize()
	lookup := squirrel.
		Select("name", "value", "rel_owner", "updated_at", "updated_by").
		From(r.dbTable).
		// Always filter by owner
		Where("rel_owner = ?", f.OwnedBy)

	if len(f.Prefix) > 0 {
		lookup = lookup.Where("name LIKE ?", f.Prefix+"%")
	}

	if f.Page > 0 {
		lookup = lookup.Offset(f.PerPage * f.Page)
	}

	if f.PerPage > 0 {
		lookup = lookup.Limit(f.PerPage)
	}

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for settings")
	} else if err = r.db().Select(&ss, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not find settings")
	} else {
		return ss, nil
	}
}

func (r *repository) Set(value *Value) error {
	value.UpdatedAt = time.Now()
	return r.db().Replace(r.dbTable, value)
}

func (r *repository) Get(name string, ownedBy uint64) (value *Value, err error) {
	lookup := squirrel.
		Select("value").
		From(r.dbTable).
		Where("rel_owner = ?", ownedBy).
		Where("name = ?", name)

	value = &Value{}

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for settings")
	} else if err = r.db().Get(value, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not get settings")
	} else {
		return value, nil
	}

}
