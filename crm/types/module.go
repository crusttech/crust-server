package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Modules - CRM module definitions
	Module struct {
		ID     uint64         `json:"moduleID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Meta   types.JSONText `json:"meta" db:"json"`
		Fields ModuleFieldSet `json:"fields" db:"-"`
		Page   *Page          `json:"page,omitempty"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// Modules - CRM module definitions
	ModuleField struct {
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`
		Place    int    `json:"-" db:"place"`

		Kind  string `json:"kind" db:"kind"`
		Name  string `json:"name" db:"name"`
		Label string `json:"label" db:"label"`

		Options types.JSONText `json:"options" db:"json"`

		Private  bool `json:"isPrivate" db:"is_private"`
		Required bool `json:"isRequired" db:"is_required"`
		Visible  bool `json:"isVisible" db:"is_visible"`
		Multi    bool `json:"isMulti" db:"is_multi"`
	}
)

func (set *ModuleFieldSet) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, set)
	}
	return nil
}

func (set ModuleFieldSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (set ModuleFieldSet) Names() (names []string) {
	names = make([]string, len(set))

	for i := range set {
		names[i] = set[i].Name
	}

	return
}

func (set ModuleFieldSet) HasName(name string) bool {
	for i := range set {
		if name == set[i].Name {
			return true
		}
	}

	return false
}

func (set ModuleFieldSet) FindByName(name string) *ModuleField {
	for i := range set {
		if name == set[i].Name {
			return set[i]
		}
	}

	return nil
}

func (set ModuleFieldSet) FilterByModule(moduleID uint64) (ff ModuleFieldSet) {
	for i := range set {
		if set[i].ModuleID == moduleID {
			ff = append(ff, set[i])
		}
	}

	return
}

// IsRef tells us if value of this field be a reference to something (another record, user)?
func (f ModuleField) IsRef() bool {
	return f.Kind == "Record" || f.Kind == "Owner" || f.Kind == "File"
}

// Resource returns a system resource ID for this type
func (r *Module) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "module",
		ID:      r.ID,
	}

	return resource
}
