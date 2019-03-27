package settings

import (
	"context"
	"encoding/json"

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

		Get(name string, value interface{}) (bool, error)
		Set(name string, value interface{}) error
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
		From(r.dbTable)

	// Always filter by owner
	lookup.Where("rel_owner = ?", f.OwnedBy)

	if len(f.Prefix) > 0 {
		lookup.Where("name LIKE ?", f.Prefix+"%")
	}

	if f.Page > 0 {
		lookup.Offset(f.PerPage * f.Page)
	}

	if f.PerPage > 0 {
		lookup.Limit(f.PerPage)
	}

	if query, args, err := lookup.ToSql(); err != nil {
		return nil, errors.Wrap(err, "could not build lookup query for settings")
	} else if err = r.db().Select(&ss, query, args...); err != nil {
		return nil, errors.Wrap(err, "could not find settings")
	} else {
		return ss, nil
	}
}

func (r *repository) Set(name string, value interface{}) error {
	if jsonValue, err := json.Marshal(value); err != nil {
		return errors.Wrap(err, "Error marshaling settings value")
	} else {
		return r.db().Replace(r.dbTable, struct {
			Key string          `db:"name"`
			Val json.RawMessage `db:"value"`
		}{name, jsonValue})
	}
}

func (r *repository) Get(name string, value interface{}) (bool, error) {
	sql := "SELECT value FROM " + r.dbTable + " WHERE name = ?"

	var stored json.RawMessage

	if err := r.db().Get(&stored, sql, name); err != nil {
		return false, errors.Wrap(err, "Error reading settings from the database")
	} else if stored == nil {
		return false, nil
	} else if err := json.Unmarshal(stored, value); err != nil {
		return false, errors.Wrap(err, "Error unmarshaling settings value")
	}

	return true, nil
}
