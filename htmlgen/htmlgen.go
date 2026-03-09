// Package htmlgen generates standalone HTML files for ERD visualization.
package htmlgen

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/grokify/entscape/schema"
)

//go:embed templates/standalone.html
var templatesFS embed.FS

// Options configures HTML generation.
type Options struct {
	// Title is the page title (defaults to "Entscape")
	Title string
	// SourceURL is an optional link to the source repository
	SourceURL string
}

// TemplateData holds data passed to the HTML template.
type TemplateData struct {
	Title      string
	SourceURL  string
	SchemaJSON template.JS
}

// Generate creates a standalone HTML file from the schema.
func Generate(s *schema.Schema, opts Options) ([]byte, error) {
	if opts.Title == "" {
		opts.Title = "Entscape"
	}

	// Marshal schema to JSON
	schemaBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("marshaling schema: %w", err)
	}

	// Parse template
	tmplContent, err := templatesFS.ReadFile("templates/standalone.html")
	if err != nil {
		return nil, fmt.Errorf("reading template: %w", err)
	}

	tmpl, err := template.New("standalone").Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	// Execute template
	data := TemplateData{
		Title:      opts.Title,
		SourceURL:  opts.SourceURL,
		SchemaJSON: template.JS(schemaBytes),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}
