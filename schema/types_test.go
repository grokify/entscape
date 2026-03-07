package schema

import (
	"encoding/json"
	"testing"
)

func TestNewSchema(t *testing.T) {
	s := NewSchema()

	if s.Version != SchemaVersion {
		t.Errorf("Version = %q, want %q", s.Version, SchemaVersion)
	}

	if s.Entities == nil {
		t.Error("Entities should not be nil")
	}

	if len(s.Entities) != 0 {
		t.Errorf("len(Entities) = %d, want 0", len(s.Entities))
	}
}

func TestSchemaJSON(t *testing.T) {
	s := &Schema{
		Version: "1",
		Package: &Package{
			Name:   "example",
			Source: "https://github.com/example/repo",
		},
		Entities: []Entity{
			{
				Name:        "User",
				Description: "User entity",
				Fields: []Field{
					{Name: "id", Type: "int"},
					{Name: "email", Type: "string", Attrs: []string{AttrUnique, AttrRequired}},
				},
				Edges: []Edge{
					{Name: "posts", Target: "Post", Relation: RelationO2M},
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Unmarshal back
	var s2 Schema
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Verify round-trip
	if s2.Version != s.Version {
		t.Errorf("Version = %q, want %q", s2.Version, s.Version)
	}
	if s2.Package.Name != s.Package.Name {
		t.Errorf("Package.Name = %q, want %q", s2.Package.Name, s.Package.Name)
	}
	if len(s2.Entities) != 1 {
		t.Fatalf("len(Entities) = %d, want 1", len(s2.Entities))
	}
	if s2.Entities[0].Name != "User" {
		t.Errorf("Entities[0].Name = %q, want %q", s2.Entities[0].Name, "User")
	}
}

func TestRelationConstants(t *testing.T) {
	if RelationO2O != "O2O" {
		t.Errorf("RelationO2O = %q, want O2O", RelationO2O)
	}
	if RelationO2M != "O2M" {
		t.Errorf("RelationO2M = %q, want O2M", RelationO2M)
	}
	if RelationM2O != "M2O" {
		t.Errorf("RelationM2O = %q, want M2O", RelationM2O)
	}
	if RelationM2M != "M2M" {
		t.Errorf("RelationM2M = %q, want M2M", RelationM2M)
	}
}

func TestAttrConstants(t *testing.T) {
	attrs := []string{
		AttrPrimary,
		AttrUnique,
		AttrRequired,
		AttrOptional,
		AttrImmutable,
		AttrSensitive,
		AttrNillable,
		AttrDefault,
	}

	// Verify all are non-empty
	for _, attr := range attrs {
		if attr == "" {
			t.Error("Attribute constant should not be empty")
		}
	}
}

func TestGenerateJSONSchema(t *testing.T) {
	data, err := GenerateJSONSchema()
	if err != nil {
		t.Fatalf("GenerateJSONSchema failed: %v", err)
	}

	// Verify it's valid JSON
	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Generated schema is not valid JSON: %v", err)
	}

	// Verify expected fields
	if schema["$schema"] == nil {
		t.Error("Missing $schema field")
	}
	if schema["$id"] == nil {
		t.Error("Missing $id field")
	}
	if schema["title"] != "Entscape Schema" {
		t.Errorf("title = %v, want Entscape Schema", schema["title"])
	}
}
