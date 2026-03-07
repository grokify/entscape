// Package schema defines the entscape JSON schema types.
package schema

// Schema is the root structure for entscape visualization data.
type Schema struct {
	Version  string   `json:"version"`
	Package  *Package `json:"package,omitempty"`
	Entities []Entity `json:"entities"`
}

// Package contains metadata about the Go package.
type Package struct {
	Name   string `json:"name"`
	Source string `json:"source,omitempty"`
	Branch string `json:"branch,omitempty"`
	Docs   string `json:"docs,omitempty"`
}

// Entity represents an Ent entity (table).
type Entity struct {
	Name        string   `json:"name"`
	Path        string   `json:"path,omitempty"`
	Description string   `json:"description,omitempty"`
	Fields      []Field  `json:"fields,omitempty"`
	Edges       []Edge   `json:"edges,omitempty"`
	Indexes     []Index  `json:"indexes,omitempty"`
	Mixins      []string `json:"mixins,omitempty"`
}

// Field represents an entity field (column).
type Field struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Attrs       []string `json:"attrs,omitempty"`
	Description string   `json:"description,omitempty"`
	Default     string   `json:"default,omitempty"`
}

// Edge represents an entity edge (relationship).
type Edge struct {
	Name        string `json:"name"`
	Target      string `json:"target"`
	Relation    string `json:"relation"`
	Inverse     string `json:"inverse,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Unique      bool   `json:"unique,omitempty"`
	Description string `json:"description,omitempty"`
}

// Index represents an entity index.
type Index struct {
	Fields []string `json:"fields"`
	Unique bool     `json:"unique,omitempty"`
	Edges  []string `json:"edges,omitempty"`
}

// Relation types for edges.
const (
	RelationO2O = "O2O" // One-to-One
	RelationO2M = "O2M" // One-to-Many
	RelationM2O = "M2O" // Many-to-One
	RelationM2M = "M2M" // Many-to-Many
)

// Field attributes.
const (
	AttrPrimary   = "primary"
	AttrUnique    = "unique"
	AttrRequired  = "required"
	AttrOptional  = "optional"
	AttrImmutable = "immutable"
	AttrSensitive = "sensitive"
	AttrNillable  = "nillable"
	AttrDefault   = "default"
)

// SchemaVersion is the current schema version.
const SchemaVersion = "1"

// NewSchema creates a new Schema with the current version.
func NewSchema() *Schema {
	return &Schema{
		Version:  SchemaVersion,
		Entities: []Entity{},
	}
}
