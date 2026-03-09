// Package export provides JSON export functionality for entscape schemas.
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/entscape/schema"
)

// Options configures the export behavior.
type Options struct {
	RepoURL string // Repository URL (e.g., https://github.com/org/repo)
	Branch  string // Git branch for source links
	DocsURL string // Documentation base URL
	Indent  bool   // Pretty-print JSON
}

// Exporter exports schemas to JSON.
type Exporter struct {
	opts Options
}

// New creates a new Exporter with the given options.
func New(opts Options) *Exporter {
	if opts.Branch == "" {
		opts.Branch = "main"
	}
	return &Exporter{opts: opts}
}

// Export converts a schema to JSON bytes.
func (e *Exporter) Export(s *schema.Schema) ([]byte, error) {
	// Add package metadata if repo URL is provided
	if e.opts.RepoURL != "" {
		if s.Package == nil {
			s.Package = &schema.Package{}
		}
		s.Package.Source = e.opts.RepoURL
		s.Package.Branch = e.opts.Branch

		if e.opts.DocsURL != "" {
			s.Package.Docs = e.opts.DocsURL
		}

		// Generate source links for each entity
		e.addSourceLinks(s)
	}

	if e.opts.Indent {
		return json.MarshalIndent(s, "", "  ")
	}
	return json.Marshal(s)
}

// ExportToFile writes the schema to a file.
func (e *Exporter) ExportToFile(s *schema.Schema, path string) error {
	data, err := e.Export(s)
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// AddSourceLinks adds source URLs to entity paths and returns the modified schema.
func (e *Exporter) AddSourceLinks(s *schema.Schema) *schema.Schema {
	// Add package metadata
	if e.opts.RepoURL != "" {
		if s.Package == nil {
			s.Package = &schema.Package{}
		}
		s.Package.Source = e.opts.RepoURL
		s.Package.Branch = e.opts.Branch
	}

	e.addSourceLinks(s)
	return s
}

// addSourceLinks adds source URLs to entity paths.
func (e *Exporter) addSourceLinks(s *schema.Schema) {
	for i := range s.Entities {
		entity := &s.Entities[i]
		if entity.Path != "" {
			// Convert path to full source URL
			// e.g., schema/user.go -> https://github.com/org/repo/blob/main/ent/schema/user.go
			entity.Path = e.buildSourceURL(entity.Path)
		}
	}
}

// buildSourceURL constructs a source file URL.
func (e *Exporter) buildSourceURL(path string) string {
	// Detect provider and build appropriate URL
	repoURL := strings.TrimSuffix(e.opts.RepoURL, "/")

	switch {
	case strings.Contains(repoURL, "github.com"):
		return fmt.Sprintf("%s/blob/%s/ent/%s", repoURL, e.opts.Branch, path)
	case strings.Contains(repoURL, "gitlab"):
		return fmt.Sprintf("%s/-/blob/%s/ent/%s", repoURL, e.opts.Branch, path)
	case strings.Contains(repoURL, "bitbucket"):
		return fmt.Sprintf("%s/src/%s/ent/%s", repoURL, e.opts.Branch, path)
	default:
		// Generic format
		return fmt.Sprintf("%s/blob/%s/ent/%s", repoURL, e.opts.Branch, path)
	}
}

// BuildDocsURL constructs a documentation URL for an entity.
func BuildDocsURL(baseURL, entityName string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	// pkg.go.dev format: https://pkg.go.dev/github.com/org/repo/ent#User
	return fmt.Sprintf("%s#%s", baseURL, entityName)
}
