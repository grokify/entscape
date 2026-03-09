# Entscape

Interactive entity-relationship diagram visualization for [Ent](https://entgo.io) schemas.

## Features

- **Parse Ent Schemas** - Extract entities, fields, edges, and relationships from Go code
- **Interactive Visualization** - Explore your schema with Cytoscape.js-powered diagrams
- **Source Links** - Click entities to jump to source code on GitHub/GitLab/Bitbucket
- **Static HTML** - Generate standalone HTML files for GitHub Pages deployment
- **JSON Export** - Export schema data for custom tooling

## Quick Start

```bash
# Install
go install github.com/grokify/entscape/cmd/entscape@latest

# Generate interactive HTML
entscape html ./ent/schema -o docs/erd.html --repo https://github.com/your/repo

# Or start a development server
entscape serve ./ent/schema --web web
```

## Example Output

See the [Demo](demo.md) page for a live example of the generated visualization.

## How It Works

1. **Parse** - Entscape reads your Ent schema Go files using AST analysis
2. **Extract** - Entities, fields, edges, indexes, and relationships are extracted
3. **Generate** - Output as JSON or standalone HTML with embedded visualization
4. **Deploy** - Single HTML file works on any static hosting (GitHub Pages, Netlify, etc.)

## Architecture

```
Ent Schema Files (Go)
        │
        ▼
   ┌─────────┐
   │  Parser │  AST-based extraction
   └────┬────┘
        │
        ▼
   ┌─────────┐
   │  Schema │  Normalized data model
   └────┬────┘
        │
        ├──────────────┬──────────────┐
        ▼              ▼              ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │  JSON   │   │  HTML   │   │   NPM   │
   │ Export  │   │Generator│   │ Module  │
   └─────────┘   └─────────┘   └─────────┘
```

## License

MIT License - see [LICENSE](https://github.com/grokify/entscape/blob/main/LICENSE)
