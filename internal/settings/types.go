package settings

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	Value struct {
		Name  string         `json:"name"  db:"name"`
		Value types.JSONText `json:"value" db:"value"`

		// Setting owner, 0 for global settings
		OwnedBy uint64 `json:"ownedBy" db:"rel_owner"`

		// Who updated & when
		UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
		UpdatedBy uint64    `json:"updatedBy" db:"updated_by"`
	}

	Filter struct {
		Prefix  string `json:"prefix"`
		OwnedBy uint64 `json:"ownedBy"`
		Page    uint64 `json:"page"`
		PerPage uint64 `json:"perPage"`
	}

	KV map[string]types.JSONText
)

const (
	settingsFilterPerPageMax = 100
)

func (f *Filter) Normalize() {
	f.Prefix = strings.TrimSpace(f.Prefix)
	if f.PerPage > settingsFilterPerPageMax {
		f.PerPage = settingsFilterPerPageMax
	}
}

func (v *Value) SetValueAsString(str string) error {
	var dummy interface{}
	// Test input to be sure we can save it...
	if err := json.Unmarshal([]byte(str), &dummy); err != nil {
		return err
	}

	v.Value = types.JSONText(str)
	return nil
}

func (ss ValueSet) KV() KV {
	m := KV{}

	_ = ss.Walk(func(v *Value) error {
		m[v.Name] = v.Value
		return nil
	})

	return m
}

func (kv KV) Bool(k string) (out bool) {
	if v, ok := kv[k]; ok {
		v.Unmarshal(&out)
	}

	return
}

func (kv KV) String(k string) (out string) {
	if v, ok := kv[k]; ok {
		v.Unmarshal(&out)
	}

	return
}

func (kv KV) Filter(prefix string) KV {
	var out = KV{}
	for k, v := range kv {
		if strings.Index(k, prefix) == 0 {
			out[k] = v
		}
	}

	return out
}