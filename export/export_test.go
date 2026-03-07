package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/grokify/entscape/schema"
)

func TestExport(t *testing.T) {
	s := &schema.Schema{
		Version: "1",
		Entities: []schema.Entity{
			{
				Name: "User",
				Fields: []schema.Field{
					{Name: "id", Type: "int"},
					{Name: "email", Type: "string", Attrs: []string{"unique"}},
				},
				Edges: []schema.Edge{
					{Name: "posts", Target: "Post", Relation: schema.RelationO2M},
				},
			},
		},
	}

	exp := New(Options{Indent: true})
	data, err := exp.Export(s)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify JSON is valid
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	// Verify structure
	if result["version"] != "1" {
		t.Errorf("version = %v, want 1", result["version"])
	}

	entities, ok := result["entities"].([]interface{})
	if !ok || len(entities) != 1 {
		t.Errorf("entities = %v, want 1 entity", result["entities"])
	}
}

func TestExportWithRepoURL(t *testing.T) {
	s := &schema.Schema{
		Version: "1",
		Entities: []schema.Entity{
			{
				Name: "User",
				Path: "schema/user.go",
			},
		},
	}

	exp := New(Options{
		RepoURL: "https://github.com/example/repo",
		Branch:  "main",
		Indent:  true,
	})

	data, err := exp.Export(s)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify source link was added
	if !strings.Contains(string(data), "https://github.com/example/repo/blob/main/ent/schema/user.go") {
		t.Errorf("Expected GitHub source link in output:\n%s", string(data))
	}

	// Verify package metadata
	if !strings.Contains(string(data), `"source": "https://github.com/example/repo"`) {
		t.Errorf("Expected package source in output:\n%s", string(data))
	}
}

func TestBuildSourceURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		branch  string
		path    string
		want    string
	}{
		{
			name:    "GitHub",
			repoURL: "https://github.com/org/repo",
			branch:  "main",
			path:    "schema/user.go",
			want:    "https://github.com/org/repo/blob/main/ent/schema/user.go",
		},
		{
			name:    "GitLab",
			repoURL: "https://gitlab.com/org/repo",
			branch:  "develop",
			path:    "schema/post.go",
			want:    "https://gitlab.com/org/repo/-/blob/develop/ent/schema/post.go",
		},
		{
			name:    "Bitbucket",
			repoURL: "https://bitbucket.org/org/repo",
			branch:  "master",
			path:    "schema/item.go",
			want:    "https://bitbucket.org/org/repo/src/master/ent/schema/item.go",
		},
		{
			name:    "Trailing slash",
			repoURL: "https://github.com/org/repo/",
			branch:  "main",
			path:    "schema/user.go",
			want:    "https://github.com/org/repo/blob/main/ent/schema/user.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := New(Options{RepoURL: tt.repoURL, Branch: tt.branch})
			got := exp.buildSourceURL(tt.path)
			if got != tt.want {
				t.Errorf("buildSourceURL(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestBuildDocsURL(t *testing.T) {
	tests := []struct {
		baseURL    string
		entityName string
		want       string
	}{
		{
			baseURL:    "https://pkg.go.dev/github.com/org/repo/ent",
			entityName: "User",
			want:       "https://pkg.go.dev/github.com/org/repo/ent#User",
		},
		{
			baseURL:    "https://pkg.go.dev/github.com/org/repo/ent/",
			entityName: "Post",
			want:       "https://pkg.go.dev/github.com/org/repo/ent#Post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.entityName, func(t *testing.T) {
			got := BuildDocsURL(tt.baseURL, tt.entityName)
			if got != tt.want {
				t.Errorf("BuildDocsURL(%q, %q) = %q, want %q", tt.baseURL, tt.entityName, got, tt.want)
			}
		})
	}
}

func TestDefaultBranch(t *testing.T) {
	exp := New(Options{})
	if exp.opts.Branch != "main" {
		t.Errorf("default branch = %q, want %q", exp.opts.Branch, "main")
	}
}
