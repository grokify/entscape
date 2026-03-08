# entscape Tasks

Interactive Ent.go schema visualization using Cytoscape.js.

## Overview

entscape provides interactive entity-relationship diagram visualization for [Ent](https://entgo.io) schemas, publishable to GitHub Pages or any static hosting.

**Architecture:**

- Go CLI reads Ent schema and exports JSON
- NPM module renders interactive Cytoscape.js graph
- Static HTML generator for quick deployment

## Phase 1: Core

### 1.1 Project Setup
**Status**: Complete

- [x] Initialize Go module (`go mod init github.com/grokify/entscape`)
- [ ] Initialize NPM package (`npm init` in `web/`)
- [x] Set up project structure
- [x] Add LICENSE (MIT)
- [x] Add README.md with project overview

### 1.2 JSON Schema Definition
**Status**: Complete

Define `entscape.schema.json` format:

- [x] Create JSON Schema for validation (`schema/entscape.schema.json`)
- [x] Define entity structure (name, path, fields, edges)
- [x] Define field structure (name, type, attrs)
- [x] Define edge structure (name, target, relation, inverse)
- [x] Define package metadata (name, source URL, branch)
- [x] Add Go types matching the schema (`schema/types.go`)
- [x] Add schema version (v1)
- [x] Add `entscape schema` CLI command to regenerate JSON Schema

### 1.3 Ent Schema Parser (Go)
**Status**: Complete

Parse Ent schema files and extract metadata:

- [x] Read ent/schema directory
- [x] Parse Go files using `go/ast`
- [x] Extract entity names from type definitions
- [x] Extract fields (name, type, modifiers)
- [x] Extract edges (name, target type, relation type)
- [x] Detect relation types (O2O, O2M, M2O, M2M)
- [x] Detect inverse edges
- [x] Handle edge.From and edge.To
- [x] Unit tests with sample schemas

### 1.4 JSON Exporter (Go)
**Status**: Complete

Export parsed schema to JSON:

- [x] Convert parsed schema to JSON types
- [x] Generate source links (GitHub/GitLab URLs)
- [x] Support `--repo`, `--branch` flags
- [x] Pretty-print JSON output
- [ ] Validate output against JSON Schema
- [x] Write to file or stdout
- [x] Add unit tests for export package

### 1.5 Cytoscape.js Renderer (TypeScript/NPM)
**Status**: Complete

NPM module for rendering:

- [x] Set up TypeScript project
- [x] Install Cytoscape.js and dagre layout
- [x] Define TypeScript types matching JSON schema
- [x] Create `entscape.render(container, schema, options)` function
- [x] Implement dagre layout for ER diagrams
- [x] Style nodes as entity boxes with field lists
- [x] Style edges with relation labels (1:1, 1:N, M:N)
- [x] Add click handlers to open source links
- [x] Export as ES module and UMD
- [ ] Publish to NPM (`@grokify/entscape`)

### 1.6 CLI Integration
**Status**: Partial

Go CLI that ties it together:

- [x] `entscape generate` - Parse and export JSON
- [ ] `entscape serve` - Local dev server with hot reload
- [x] Flags: `--schema`, `--output`, `--repo`, `--branch`
- [x] Error handling and validation
- [ ] `--html` flag for standalone HTML output

## Phase 2: Static HTML & Docs Integration

### 2.1 Static HTML Generator
**Status**: Not Started

Generate standalone HTML file:

- [ ] Create HTML template with embedded Cytoscape.js
- [ ] Embed JSON schema inline or as separate file
- [ ] Bundle NPM module into single JS file
- [ ] `entscape generate --html` outputs complete HTML
- [ ] Support custom title, theme options
- [ ] Optimize for GitHub Pages deployment

### 2.2 Source Link Integration
**Status**: Partial

Link to source code repositories:

- [ ] Auto-detect GitHub/GitLab from go.mod
- [x] Generate correct blob URLs (GitHub, GitLab, Bitbucket)
- [ ] Support custom URL templates
- [ ] Open source file on entity click
- [ ] Highlight specific line if possible

### 2.3 Documentation Link Integration
**Status**: Not Started

Link to Go documentation:

- [ ] Generate pkg.go.dev links for public packages
- [ ] Support custom docs URL (`--docs-base`) for private repos
- [ ] Add docs link to entity tooltip/panel
- [ ] Support pkgsite self-hosted instances

### 2.4 GitHub Pages Deployment
**Status**: Not Started

Easy deployment workflow:

- [ ] Generate to `docs/` folder by default
- [ ] GitHub Action for auto-regeneration on schema changes
- [ ] Example workflow YAML
- [ ] Documentation for setup

## Phase 3: Enhanced Features

### 3.1 Interactive Features
**Status**: Not Started

Improve visualization UX:

- [ ] Search/filter entities by name
- [ ] Zoom to fit / reset view button
- [ ] Minimap navigator for large schemas
- [ ] Collapse/expand entity fields
- [ ] Highlight connected entities on hover
- [ ] Keyboard navigation

### 3.2 Entity Details Panel
**Status**: Not Started

Show detailed info on selection:

- [ ] Side panel with full entity details
- [ ] List all fields with types and attributes
- [ ] List all edges with targets
- [ ] Show indexes, constraints
- [ ] Copy field/edge definitions

### 3.3 Export Options
**Status**: Not Started

Export visualization:

- [ ] Export as PNG image
- [ ] Export as SVG
- [ ] Export as Mermaid diagram
- [ ] Export as D2 diagram
- [ ] Copy to clipboard

### 3.4 Theming
**Status**: Not Started

Visual customization:

- [ ] Light/dark mode toggle
- [ ] Custom color schemes
- [ ] Entity grouping by package/tag
- [ ] Highlight specific entities via URL params

### 3.5 Advanced Schema Support
**Status**: Not Started

Handle complex Ent features:

- [ ] Mixin fields and edges
- [ ] Privacy policies (visual indicator)
- [ ] Hooks (visual indicator)
- [ ] Annotations
- [ ] Multi-schema projects

## Project Structure

```
entscape/
├── cmd/entscape/           # Go CLI
│   └── main.go
├── parser/                 # Ent schema parser
│   ├── parser.go
│   └── parser_test.go
├── schema/                 # JSON schema types
│   ├── types.go
│   ├── entscape.schema.json
│   └── validate.go
├── export/                 # JSON exporter
│   ├── export.go
│   └── links.go
├── web/                    # NPM module
│   ├── src/
│   │   ├── index.ts
│   │   ├── types.ts
│   │   ├── render.ts
│   │   └── styles.ts
│   ├── package.json
│   ├── tsconfig.json
│   └── rollup.config.js
├── templates/              # HTML templates
│   └── standalone.html
├── testdata/               # Test Ent schemas
│   └── basic/
│       └── schema/
├── go.mod
├── go.sum
├── TASKS.md
├── README.md
└── LICENSE
```

## Dependencies

### Go
- `golang.org/x/tools/go/packages` - Go package loading
- `entgo.io/ent` - Ent types (optional, for validation)

### NPM
- `cytoscape` - Graph visualization
- `cytoscape-dagre` - Hierarchical layout
- `typescript` - Type safety
- `rollup` - Bundling

## Implementation Order

1. Project setup and JSON schema definition
2. Go parser for Ent schemas
3. JSON exporter with source links
4. NPM module with Cytoscape.js
5. CLI integration
6. Static HTML generator
7. Docs links and GitHub Pages
8. Enhanced features

## Notes

- Prioritize clean, documented JSON schema that others can use
- NPM module should work standalone without Go CLI
- Support both ES modules and script tag inclusion
- Test with real-world Ent schemas of varying complexity
