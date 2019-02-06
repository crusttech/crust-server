package types

import (
	"time"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Organisations - Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.
	Organisation struct {
		ID         uint64     `json:"organisationID,string" db:"id"`
		FQN        string     `json:"fqn" db:"fqn"`
		Name       string     `json:"name" db:"name"`
		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	OrganisationFilter struct {
		Query string
	}
)

// Resource returns a system resource ID for this type
func (r *Organisation) Resource() rules.Resource {
	resource := rules.Resource{
		Scope: "organisation",
	}
	if r != nil {
		resource.ID = r.ID
		resource.Name = r.Name
	}
	return resource
}
