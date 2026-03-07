// Package parser parses Ent schema files using Go AST.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/entscape/schema"
)

// Parser parses Ent schema directories.
type Parser struct {
	fset *token.FileSet
}

// New creates a new Parser.
func New() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseDir parses all Ent schema files in a directory.
func (p *Parser) ParseDir(dir string) (*schema.Schema, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	s := schema.NewSchema()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		entity, err := p.parseFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		if entity != nil {
			entity.Path = filepath.Join("schema", entry.Name())
			s.Entities = append(s.Entities, *entity)
		}
	}

	// Resolve relation types based on edge analysis
	p.resolveRelations(s)

	return s, nil
}

// parseFile parses a single Ent schema file.
func (p *Parser) parseFile(path string) (*schema.Entity, error) {
	f, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	var entity *schema.Entity

	// Find the schema struct type
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Check if it embeds ent.Schema
			if !hasEntSchemaEmbed(structType) {
				continue
			}

			entity = &schema.Entity{
				Name: typeSpec.Name.Name,
			}

			// Extract description from comments
			if genDecl.Doc != nil {
				entity.Description = cleanComment(genDecl.Doc.Text())
			}

			break
		}
	}

	if entity == nil {
		return nil, nil // Not an Ent schema file
	}

	// Find Fields() and Edges() methods
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if it's a method on our entity type
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}

		switch funcDecl.Name.Name {
		case "Fields":
			entity.Fields = p.parseFields(funcDecl)
		case "Edges":
			entity.Edges = p.parseEdges(funcDecl)
		case "Indexes":
			entity.Indexes = p.parseIndexes(funcDecl)
		case "Mixin":
			entity.Mixins = p.parseMixins(funcDecl)
		}
	}

	return entity, nil
}

// parseFields extracts fields from a Fields() method.
func (p *Parser) parseFields(fn *ast.FuncDecl) []schema.Field {
	var fields []schema.Field

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if this call chain contains a field.X() call at its base
		baseCall, fieldType := p.findFieldBaseCall(call)
		if baseCall == nil {
			return true
		}

		field := p.parseFieldCall(call, baseCall, fieldType)
		if field != nil {
			fields = append(fields, *field)
		}

		return false // Don't recurse into field chain
	})

	return fields
}

// findFieldBaseCall walks down a call chain to find field.X() at its base.
// Returns the base call and the field type (e.g., "String", "Int").
func (p *Parser) findFieldBaseCall(call *ast.CallExpr) (*ast.CallExpr, string) {
	current := call
	for {
		sel, ok := current.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, ""
		}

		// Check if we've reached field.X()
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "field" {
			return current, sel.Sel.Name
		}

		// Go deeper into the chain
		innerCall, ok := sel.X.(*ast.CallExpr)
		if !ok {
			return nil, ""
		}
		current = innerCall
	}
}

// parseFieldCall parses a field call chain, extracting name, type, and modifiers.
// outerCall is the outermost call in the chain, baseCall is the field.X() call.
func (p *Parser) parseFieldCall(outerCall, baseCall *ast.CallExpr, fieldType string) *schema.Field {
	if len(baseCall.Args) == 0 {
		return nil
	}

	// Get field name from first argument of the base call
	nameLit, ok := baseCall.Args[0].(*ast.BasicLit)
	if !ok || nameLit.Kind != token.STRING {
		return nil
	}

	name := strings.Trim(nameLit.Value, `"`)

	field := &schema.Field{
		Name: name,
		Type: fieldTypeToGoType(fieldType),
	}

	// Extract modifiers from the entire chain (outer to inner)
	p.extractFieldModifiers(outerCall, field)

	return field
}

// extractFieldModifiers walks the call chain from outer to inner to extract modifiers.
func (p *Parser) extractFieldModifiers(expr ast.Expr, field *schema.Field) {
	for {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// Stop when we reach field.X()
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "field" {
			return
		}

		switch sel.Sel.Name {
		case "Unique":
			field.Attrs = append(field.Attrs, schema.AttrUnique)
		case "NotEmpty", "Required":
			field.Attrs = append(field.Attrs, schema.AttrRequired)
		case "Optional":
			field.Attrs = append(field.Attrs, schema.AttrOptional)
		case "Immutable":
			field.Attrs = append(field.Attrs, schema.AttrImmutable)
		case "Sensitive":
			field.Attrs = append(field.Attrs, schema.AttrSensitive)
		case "Nillable":
			field.Attrs = append(field.Attrs, schema.AttrNillable)
		case "Default":
			field.Attrs = append(field.Attrs, schema.AttrDefault)
		}

		// Continue down the chain
		expr = sel.X
	}
}

// parseEdges extracts edges from an Edges() method.
func (p *Parser) parseEdges(fn *ast.FuncDecl) []schema.Edge {
	var edges []schema.Edge

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check if this call chain contains an edge.To() or edge.From() call at its base
		baseCall, edgeDir := p.findEdgeBaseCall(call)
		if baseCall == nil {
			return true
		}

		edge := p.parseEdgeCall(call, baseCall, edgeDir)
		if edge != nil {
			edges = append(edges, *edge)
		}

		return false
	})

	return edges
}

// findEdgeBaseCall walks down a call chain to find edge.To() or edge.From() at its base.
// Returns the base call and the edge direction ("To" or "From").
func (p *Parser) findEdgeBaseCall(call *ast.CallExpr) (*ast.CallExpr, string) {
	current := call
	for {
		sel, ok := current.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, ""
		}

		// Check if we've reached edge.To() or edge.From()
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "edge" {
			if sel.Sel.Name == "To" || sel.Sel.Name == "From" {
				return current, sel.Sel.Name
			}
			return nil, ""
		}

		// Go deeper into the chain
		innerCall, ok := sel.X.(*ast.CallExpr)
		if !ok {
			return nil, ""
		}
		current = innerCall
	}
}

// parseEdgeCall parses an edge call chain, extracting name, target, and modifiers.
// outerCall is the outermost call in the chain, baseCall is the edge.To() or edge.From() call.
func (p *Parser) parseEdgeCall(outerCall, baseCall *ast.CallExpr, edgeDir string) *schema.Edge {
	if len(baseCall.Args) < 2 {
		return nil
	}

	// Get edge name from first argument of the base call
	nameLit, ok := baseCall.Args[0].(*ast.BasicLit)
	if !ok || nameLit.Kind != token.STRING {
		return nil
	}
	name := strings.Trim(nameLit.Value, `"`)

	// Get target type from second argument (e.g., Post.Type)
	target := extractTypeName(baseCall.Args[1])
	if target == "" {
		return nil
	}

	edge := &schema.Edge{
		Name:   name,
		Target: target,
	}

	// Determine initial relation based on edge direction
	if edgeDir == "To" {
		edge.Relation = schema.RelationO2M // Default, may be refined
	} else {
		edge.Relation = schema.RelationM2O // edge.From is typically M2O
	}

	// Extract modifiers from the entire chain (outer to inner)
	p.extractEdgeModifiers(outerCall, edge)

	return edge
}

// extractEdgeModifiers walks the call chain from outer to inner to extract modifiers.
func (p *Parser) extractEdgeModifiers(expr ast.Expr, edge *schema.Edge) {
	for {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// Stop when we reach edge.To() or edge.From()
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "edge" {
			return
		}

		switch sel.Sel.Name {
		case "Unique":
			edge.Unique = true
			// Unique edge.To is O2O, unique edge.From is still M2O (but unique on this side)
			if edge.Relation == schema.RelationO2M {
				edge.Relation = schema.RelationO2O
			}
		case "Required":
			edge.Required = true
		case "Ref":
			if len(call.Args) > 0 {
				if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					edge.Inverse = strings.Trim(lit.Value, `"`)
				}
			}
		}

		// Continue down the chain
		expr = sel.X
	}
}

// parseIndexes extracts indexes from an Indexes() method.
func (p *Parser) parseIndexes(fn *ast.FuncDecl) []schema.Index {
	var indexes []schema.Index

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "index" {
			return true
		}

		if sel.Sel.Name == "Fields" {
			idx := p.parseIndexCall(call)
			if idx != nil {
				indexes = append(indexes, *idx)
			}
		}

		return false
	})

	return indexes
}

// parseIndexCall parses a single index.Fields() call chain.
func (p *Parser) parseIndexCall(call *ast.CallExpr) *schema.Index {
	idx := &schema.Index{}

	// Extract field names from arguments
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			idx.Fields = append(idx.Fields, strings.Trim(lit.Value, `"`))
		}
	}

	if len(idx.Fields) == 0 {
		return nil
	}

	// Walk up the call chain to find modifiers
	p.extractIndexModifiers(call, idx)

	return idx
}

// extractIndexModifiers walks the call chain to extract modifiers.
func (p *Parser) extractIndexModifiers(expr ast.Expr, idx *schema.Index) {
	for {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		switch sel.Sel.Name {
		case "Unique":
			idx.Unique = true
		case "Edges":
			for _, arg := range call.Args {
				if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					idx.Edges = append(idx.Edges, strings.Trim(lit.Value, `"`))
				}
			}
		}

		expr = sel.X
	}
}

// parseMixins extracts mixin names from a Mixin() method.
func (p *Parser) parseMixins(fn *ast.FuncDecl) []string {
	var mixins []string

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		for _, elt := range comp.Elts {
			if unary, ok := elt.(*ast.UnaryExpr); ok {
				if comp2, ok := unary.X.(*ast.CompositeLit); ok {
					if ident, ok := comp2.Type.(*ast.Ident); ok {
						mixins = append(mixins, ident.Name)
					}
				}
			}
		}

		return true
	})

	return mixins
}

// resolveRelations refines relation types based on cross-entity analysis.
func (p *Parser) resolveRelations(s *schema.Schema) {
	// Build edge lookup
	edgeMap := make(map[string]map[string]*schema.Edge)
	for i := range s.Entities {
		entity := &s.Entities[i]
		edgeMap[entity.Name] = make(map[string]*schema.Edge)
		for j := range entity.Edges {
			edge := &entity.Edges[j]
			edgeMap[entity.Name][edge.Name] = edge
		}
	}

	// Detect M2M relationships (both sides have edges to each other without Unique)
	for i := range s.Entities {
		entity := &s.Entities[i]
		for j := range entity.Edges {
			edge := &entity.Edges[j]

			if edge.Inverse != "" {
				// Check if inverse exists on target
				if targetEdges, ok := edgeMap[edge.Target]; ok {
					if inverseEdge, ok := targetEdges[edge.Inverse]; ok {
						// Both sides have edges - could be M2M
						if !edge.Unique && !inverseEdge.Unique {
							edge.Relation = schema.RelationM2M
							inverseEdge.Relation = schema.RelationM2M
						}
					}
				}
			}
		}
	}
}

// Helper functions

func hasEntSchemaEmbed(st *ast.StructType) bool {
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 { // Embedded field
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == "ent" && sel.Sel.Name == "Schema" {
						return true
					}
				}
			}
		}
	}
	return false
}

func extractTypeName(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	if sel.Sel.Name != "Type" {
		return ""
	}
	if ident, ok := sel.X.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

func fieldTypeToGoType(fieldType string) string {
	switch fieldType {
	case "Int", "Int8", "Int16", "Int32", "Int64":
		return strings.ToLower(fieldType)
	case "Uint", "Uint8", "Uint16", "Uint32", "Uint64":
		return strings.ToLower(fieldType)
	case "Float", "Float32":
		return "float32"
	case "Float64":
		return "float64"
	case "String", "Text":
		return "string"
	case "Bool":
		return "bool"
	case "Time":
		return "time.Time"
	case "UUID":
		return "uuid.UUID"
	case "Bytes":
		return "[]byte"
	case "JSON":
		return "json.RawMessage"
	case "Enum":
		return "enum"
	case "Other":
		return "any"
	default:
		return fieldType
	}
}

func cleanComment(s string) string {
	s = strings.TrimSpace(s)
	// Remove "TypeName " prefix from comments like "User holds the schema..."
	if idx := strings.Index(s, " "); idx > 0 {
		first := s[:idx]
		if len(first) > 0 && first[0] >= 'A' && first[0] <= 'Z' {
			rest := strings.TrimSpace(s[idx+1:])
			if strings.HasPrefix(rest, "holds ") || strings.HasPrefix(rest, "is ") {
				return s
			}
		}
	}
	return s
}
