package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/grokify/entscape/export"
	"github.com/grokify/entscape/htmlgen"
	"github.com/grokify/entscape/parser"
	"github.com/grokify/entscape/schema"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "entscape",
	Short: "Interactive Ent.go schema visualization",
	Long: `entscape provides interactive entity-relationship diagram visualization
for Ent schemas, publishable to GitHub Pages or any static hosting.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate [schema-dir]",
	Short: "Generate JSON from Ent schema",
	Long: `Parse Ent schema files and generate JSON output for visualization.

The schema-dir argument should point to an ent/schema directory containing
Go files with Ent entity definitions.`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

var (
	outputFlag string
	repoFlag   string
	branchFlag string
	docsFlag   string
	prettyFlag bool
)

var (
	servePortFlag   int
	serveRepoFlag   string
	serveBranchFlag string
	serveDocsFlag   string
	serveWebFlag    string
)

var (
	htmlOutputFlag string
	htmlRepoFlag   string
	htmlBranchFlag string
	htmlTitleFlag  string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(htmlCmd)
	rootCmd.AddCommand(versionCmd)

	generateCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file (default: stdout)")
	generateCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository URL for source links")
	generateCmd.Flags().StringVar(&branchFlag, "branch", "main", "Git branch for source links")
	generateCmd.Flags().StringVar(&docsFlag, "docs", "", "Documentation base URL")
	generateCmd.Flags().BoolVar(&prettyFlag, "pretty", true, "Pretty-print JSON output")

	schemaCmd.Flags().StringVarP(&schemaOutputFlag, "output", "o", "", "Output file (default: stdout)")

	serveCmd.Flags().IntVarP(&servePortFlag, "port", "p", 8080, "Port to serve on")
	serveCmd.Flags().StringVar(&serveRepoFlag, "repo", "", "Repository URL for source links")
	serveCmd.Flags().StringVar(&serveBranchFlag, "branch", "main", "Git branch for source links")
	serveCmd.Flags().StringVar(&serveDocsFlag, "docs", "", "Documentation base URL")
	serveCmd.Flags().StringVar(&serveWebFlag, "web", "", "Web directory to serve static files from")

	htmlCmd.Flags().StringVarP(&htmlOutputFlag, "output", "o", "", "Output file (default: stdout)")
	htmlCmd.Flags().StringVar(&htmlRepoFlag, "repo", "", "Repository URL for source links")
	htmlCmd.Flags().StringVar(&htmlBranchFlag, "branch", "main", "Git branch for source links")
	htmlCmd.Flags().StringVar(&htmlTitleFlag, "title", "", "Page title (default: package name or 'Entscape')")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("entscape version %s\n", version)
	},
}

var schemaOutputFlag string

var htmlCmd = &cobra.Command{
	Use:   "html [schema-dir]",
	Short: "Generate standalone HTML visualization",
	Long: `Generate a standalone HTML file with embedded schema data for visualization.

The generated HTML file can be opened directly in a browser or deployed to
GitHub Pages or any static hosting service.

Example:
  entscape html ./ent/schema --output index.html --repo https://github.com/org/repo`,
	Args: cobra.ExactArgs(1),
	RunE: runHtml,
}

func runHtml(cmd *cobra.Command, args []string) error {
	schemaDir := args[0]

	// Parse the Ent schema directory
	p := parser.New()
	s, err := p.ParseDir(schemaDir)
	if err != nil {
		return fmt.Errorf("parsing schema: %w", err)
	}

	// Add source links if repo is specified
	if htmlRepoFlag != "" {
		exp := export.New(export.Options{
			RepoURL: htmlRepoFlag,
			Branch:  htmlBranchFlag,
		})
		s = exp.AddSourceLinks(s)
	}

	// Determine title
	title := htmlTitleFlag
	if title == "" && s.Package != nil && s.Package.Name != "" {
		title = s.Package.Name
	}

	// Generate HTML
	html, err := htmlgen.Generate(s, htmlgen.Options{
		Title:     title,
		SourceURL: htmlRepoFlag,
	})
	if err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	// Write output
	if htmlOutputFlag != "" {
		if err := os.WriteFile(htmlOutputFlag, html, 0644); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote HTML to %s (%d entities)\n", htmlOutputFlag, len(s.Entities))
	} else {
		fmt.Print(string(html))
	}

	return nil
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate JSON Schema for entscape format",
	Long:  `Generate a JSON Schema file that describes the entscape JSON format.`,
	RunE:  runSchema,
}

func runSchema(cmd *cobra.Command, args []string) error {
	data, err := schema.GenerateJSONSchema()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	if schemaOutputFlag != "" {
		if err := os.WriteFile(schemaOutputFlag, data, 0644); err != nil {
			return fmt.Errorf("writing schema: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote JSON Schema to %s\n", schemaOutputFlag)
	} else {
		fmt.Println(string(data))
	}

	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve [schema-dir]",
	Short: "Start a local development server",
	Long: `Start a local HTTP server that serves the parsed schema as JSON.

The server provides:
  /api/schema.json - The parsed Ent schema as JSON
  /api/jsonschema  - The JSON Schema definition

The schema is re-parsed on each request for development convenience.`,
	Args: cobra.ExactArgs(1),
	RunE: runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
	schemaDir := args[0]

	// Verify schema directory exists
	absDir, err := filepath.Abs(schemaDir)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return fmt.Errorf("schema directory not found: %s", absDir)
	}

	// Create export options
	expOpts := export.Options{
		RepoURL: serveRepoFlag,
		Branch:  serveBranchFlag,
		DocsURL: serveDocsFlag,
		Indent:  true,
	}

	// API endpoint for schema
	http.HandleFunc("/api/schema.json", func(w http.ResponseWriter, r *http.Request) {
		p := parser.New()
		s, err := p.ParseDir(absDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("parsing schema: %v", err), http.StatusInternalServerError)
			return
		}

		exp := export.New(expOpts)
		data, err := exp.Export(s)
		if err != nil {
			http.Error(w, fmt.Sprintf("exporting schema: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_, _ = w.Write(data)
	})

	// API endpoint for JSON Schema
	http.HandleFunc("/api/jsonschema", func(w http.ResponseWriter, r *http.Request) {
		data, err := schema.GenerateJSONSchema()
		if err != nil {
			http.Error(w, fmt.Sprintf("generating schema: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_, _ = w.Write(data)
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Serve static files from web directory if specified
	if serveWebFlag != "" {
		webDir, err := filepath.Abs(serveWebFlag)
		if err != nil {
			return fmt.Errorf("resolving web path: %w", err)
		}
		if _, err := os.Stat(webDir); os.IsNotExist(err) {
			return fmt.Errorf("web directory not found: %s", webDir)
		}
		fs := http.FileServer(http.Dir(webDir))
		http.Handle("/", fs)
	}

	addr := fmt.Sprintf(":%d", servePortFlag)
	fmt.Printf("Starting server at http://localhost%s\n", addr)
	fmt.Printf("  Schema:      http://localhost%s/api/schema.json\n", addr)
	fmt.Printf("  JSON Schema: http://localhost%s/api/jsonschema\n", addr)
	fmt.Printf("  Health:      http://localhost%s/health\n", addr)
	if serveWebFlag != "" {
		fmt.Printf("  Web UI:      http://localhost%s/example.html\n", addr)
	}
	fmt.Printf("\nServing schema from: %s\n", absDir)
	fmt.Println("Press Ctrl+C to stop")

	return http.ListenAndServe(addr, nil)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	schemaDir := args[0]

	// Parse the Ent schema directory
	p := parser.New()
	schema, err := p.ParseDir(schemaDir)
	if err != nil {
		return fmt.Errorf("parsing schema: %w", err)
	}

	// Create exporter with options
	exp := export.New(export.Options{
		RepoURL: repoFlag,
		Branch:  branchFlag,
		DocsURL: docsFlag,
		Indent:  prettyFlag,
	})

	// Export to JSON
	data, err := exp.Export(schema)
	if err != nil {
		return fmt.Errorf("exporting schema: %w", err)
	}

	// Write output
	if outputFlag != "" {
		if err := os.WriteFile(outputFlag, data, 0644); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote %d entities to %s\n", len(schema.Entities), outputFlag)
	} else {
		// Pretty print to stdout
		if prettyFlag {
			var prettyJSON map[string]interface{}
			if err := json.Unmarshal(data, &prettyJSON); err != nil {
				return fmt.Errorf("formatting output: %w", err)
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(prettyJSON); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}
		} else {
			fmt.Println(string(data))
		}
	}

	return nil
}
