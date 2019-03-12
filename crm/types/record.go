package types

import (
	"time"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Record is a stored row in the `record` table
	Record struct {
		ID       uint64 `json:"recordID,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`

		Values RecordValueSet `json:"values,omitempty" db:"-"`

		OwnedBy   uint64     `db:"owned_by"   json:"ownedBy,string"`
		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}

	// RecordValue is a stored row in the `record_value` table
	RecordValue struct {
		RecordID  uint64     `db:"record_id"  json:"-"`
		Name      string     `db:"name"       json:"name"`
		Value     string     `db:"value"      json:"value,omitempty"`
		Ref       uint64     `db:"ref"        json:"-"`
		Place     uint       `db:"place"      json:"-"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}
)

// UserIDs returns a slice of user IDs from all items in the set
func (set RecordSet) UserIDs() (IDs []uint64) {
	IDs = make([]uint64, 0)

loop:
	for i := range set {
		for _, id := range IDs {
			if id == set[i].OwnedBy {
				continue loop
			}
		}

		IDs = append(IDs, set[i].OwnedBy)
	}

	return
}

func (set RecordValueSet) FilterByName(name string) (vv RecordValueSet) {
	for i := range set {
		if set[i].Name == name {
			vv = append(vv, set[i])
		}
	}

	return
}

func (set RecordValueSet) FilterByRecordID(recordID uint64) (vv RecordValueSet) {
	for i := range set {
		if set[i].RecordID == recordID {
			vv = append(vv, set[i])
		}
	}

	return
}

// Resource returns a system resource ID for this type
func (r *Record) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "module", // intentionally using module here so we can use Record's resource
		ID:      r.ModuleID,
	}

	return resource
}
