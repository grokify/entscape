package parser

import (
	"path/filepath"
	"testing"

	"github.com/grokify/entscape/schema"
)

func TestParseDir(t *testing.T) {
	p := New()
	s, err := p.ParseDir(filepath.Join("..", "testdata", "basic", "schema"))
	if err != nil {
		t.Fatalf("ParseDir failed: %v", err)
	}

	if s.Version != schema.SchemaVersion {
		t.Errorf("Version = %q, want %q", s.Version, schema.SchemaVersion)
	}

	// Should have 5 entities: User, Post, Profile, Comment, Group
	if len(s.Entities) != 5 {
		t.Errorf("len(Entities) = %d, want 5", len(s.Entities))
	}

	// Build entity map for easier testing
	entities := make(map[string]*schema.Entity)
	for i := range s.Entities {
		entities[s.Entities[i].Name] = &s.Entities[i]
	}

	// Test User entity
	user, ok := entities["User"]
	if !ok {
		t.Fatal("User entity not found")
	}

	// User should have 4 fields
	if len(user.Fields) != 4 {
		t.Errorf("User.Fields = %d, want 4", len(user.Fields))
	}

	// Check email field
	var emailField *schema.Field
	for i := range user.Fields {
		if user.Fields[i].Name == "email" {
			emailField = &user.Fields[i]
			break
		}
	}
	if emailField == nil {
		t.Fatal("email field not found")
	}
	if emailField.Type != "string" {
		t.Errorf("email.Type = %q, want %q", emailField.Type, "string")
	}
	if !containsAttr(emailField.Attrs, schema.AttrUnique) {
		t.Error("email should have unique attr")
	}

	// User should have 3 edges
	if len(user.Edges) != 3 {
		t.Errorf("User.Edges = %d, want 3", len(user.Edges))
	}

	// Check posts edge
	var postsEdge *schema.Edge
	for i := range user.Edges {
		if user.Edges[i].Name == "posts" {
			postsEdge = &user.Edges[i]
			break
		}
	}
	if postsEdge == nil {
		t.Fatal("posts edge not found")
	}
	if postsEdge.Target != "Post" {
		t.Errorf("posts.Target = %q, want %q", postsEdge.Target, "Post")
	}
	if postsEdge.Relation != schema.RelationO2M {
		t.Errorf("posts.Relation = %q, want %q", postsEdge.Relation, schema.RelationO2M)
	}

	// Check profile edge (O2O)
	var profileEdge *schema.Edge
	for i := range user.Edges {
		if user.Edges[i].Name == "profile" {
			profileEdge = &user.Edges[i]
			break
		}
	}
	if profileEdge == nil {
		t.Fatal("profile edge not found")
	}
	if profileEdge.Relation != schema.RelationO2O {
		t.Errorf("profile.Relation = %q, want %q", profileEdge.Relation, schema.RelationO2O)
	}

	// Test Post entity
	post, ok := entities["Post"]
	if !ok {
		t.Fatal("Post entity not found")
	}

	// Check author edge (M2O with inverse)
	var authorEdge *schema.Edge
	for i := range post.Edges {
		if post.Edges[i].Name == "author" {
			authorEdge = &post.Edges[i]
			break
		}
	}
	if authorEdge == nil {
		t.Fatal("author edge not found")
	}
	if authorEdge.Target != "User" {
		t.Errorf("author.Target = %q, want %q", authorEdge.Target, "User")
	}
	if authorEdge.Inverse != "posts" {
		t.Errorf("author.Inverse = %q, want %q", authorEdge.Inverse, "posts")
	}
	if !authorEdge.Required {
		t.Error("author should be required")
	}
}

func TestParseFile_NoEntSchema(t *testing.T) {
	// A regular Go file without ent.Schema should return nil
	p := New()

	// Parse a Go file that doesn't contain an Ent schema (e.g., the parser itself)
	entity, err := p.parseFile("parser.go")
	if err != nil {
		t.Fatalf("parseFile failed: %v", err)
	}
	if entity != nil {
		t.Error("expected nil entity for non-schema file, got non-nil")
	}
}

func containsAttr(attrs []string, target string) bool {
	for _, a := range attrs {
		if a == target {
			return true
		}
	}
	return false
}
