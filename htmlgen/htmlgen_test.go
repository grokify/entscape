package htmlgen

import (
	"strings"
	"testing"

	"github.com/grokify/entscape/schema"
)

func TestGenerate(t *testing.T) {
	s := &schema.Schema{
		Version: "1",
		Package: &schema.Package{
			Name:   "github.com/example/app/ent",
			Source: "https://github.com/example/app",
		},
		Entities: []schema.Entity{
			{
				Name: "User",
				Fields: []schema.Field{
					{Name: "id", Type: "int", Attrs: []string{"primary"}},
					{Name: "name", Type: "string", Attrs: []string{"required"}},
				},
				Edges: []schema.Edge{
					{Name: "posts", Target: "Post", Relation: schema.RelationO2M},
				},
			},
			{
				Name: "Post",
				Fields: []schema.Field{
					{Name: "id", Type: "int", Attrs: []string{"primary"}},
					{Name: "title", Type: "string", Attrs: []string{"required"}},
				},
			},
		},
	}

	html, err := Generate(s, Options{
		Title:     "Test App",
		SourceURL: "https://github.com/example/app",
	})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	content := string(html)

	// Check title
	if !strings.Contains(content, "<title>Test App</title>") {
		t.Error("expected title in HTML")
	}

	// Check source link
	if !strings.Contains(content, "https://github.com/example/app") {
		t.Error("expected source URL in HTML")
	}

	// Check schema data is embedded
	if !strings.Contains(content, `"name":"User"`) {
		t.Error("expected User entity in schema JSON")
	}
	if !strings.Contains(content, `"name":"Post"`) {
		t.Error("expected Post entity in schema JSON")
	}

	// Check cytoscape CDN links
	if !strings.Contains(content, "unpkg.com/cytoscape") {
		t.Error("expected cytoscape CDN link")
	}
}

func TestGenerateDefaultTitle(t *testing.T) {
	s := &schema.Schema{
		Version:  "1",
		Entities: []schema.Entity{},
	}

	html, err := Generate(s, Options{})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !strings.Contains(string(html), "<title>Entscape</title>") {
		t.Error("expected default title 'Entscape'")
	}
}
