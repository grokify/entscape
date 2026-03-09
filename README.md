# entscape

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

Interactive [Ent](https://entgo.io) schema visualization using [Cytoscape.js](https://js.cytoscape.org).

## Features

- Parse Ent schema files and export as JSON
- Interactive entity-relationship diagram with Cytoscape.js
- Click entities to open source code (GitHub/GitLab)
- Link to documentation (pkg.go.dev or self-hosted)
- Deploy to GitHub Pages as static HTML
- NPM module for embedding in existing sites

## Installation

```bash
go install github.com/grokify/entscape/cmd/entscape@latest
```

## Quick Start

```bash
# Generate JSON from Ent schema
entscape generate --schema ./ent/schema --output schema.json

# Generate static HTML
entscape generate --schema ./ent/schema --html --output docs/index.html

# With source links
entscape generate --schema ./ent/schema \
  --repo https://github.com/org/repo \
  --branch main \
  --html --output docs/index.html
```

## NPM Module

```bash
npm install @grokify/entscape
```

```javascript
import { render } from '@grokify/entscape';

const schema = await fetch('schema.json').then(r => r.json());
render(document.getElementById('diagram'), schema, {
  layout: 'dagre',
  theme: 'light'
});
```

## JSON Schema

entscape uses a documented JSON format for schema data:

```json
{
  "version": "1",
  "package": {
    "name": "github.com/org/repo/ent",
    "source": "https://github.com/org/repo",
    "branch": "main"
  },
  "entities": [
    {
      "name": "User",
      "path": "schema/user.go",
      "fields": [
        {"name": "id", "type": "int", "attrs": ["primary"]},
        {"name": "email", "type": "string", "attrs": ["unique", "required"]}
      ],
      "edges": [
        {"name": "posts", "target": "Post", "relation": "O2M", "inverse": "author"}
      ]
    }
  ]
}
```

See [schema/entscape.schema.json](schema/entscape.schema.json) for the full JSON Schema.

## GitHub Pages

1. Generate HTML to `docs/` folder:
   ```bash
   entscape generate --schema ./ent/schema --html --output docs/index.html
   ```

2. Enable GitHub Pages in repository settings (source: `docs/` folder)

3. Optionally, add GitHub Action for auto-regeneration:
   ```yaml
   name: Update Schema Diagram
   on:
     push:
       paths:
         - 'ent/schema/**'
   jobs:
     generate:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: actions/setup-go@v5
           with:
             go-version: '1.22'
         - run: go install github.com/grokify/entscape/cmd/entscape@latest
         - run: entscape generate --schema ./ent/schema --html --output docs/index.html
         - uses: stefanzweifel/git-auto-commit-action@v5
           with:
             commit_message: "docs: update schema diagram"
   ```

## Relation Types

| Code | Meaning | Example |
|------|---------|---------|
| `O2O` | One-to-One | User - Profile |
| `O2M` | One-to-Many | User - Posts |
| `M2O` | Many-to-One | Post - Author |
| `M2M` | Many-to-Many | User - Groups |

## License

MIT

 [go-ci-svg]: https://github.com/grokify/entscape/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/entscape/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/entscape/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/entscape/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/entscape/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/entscape/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/entscape
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/entscape
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/entscape
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/entscape
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fentscape
 [loc-svg]: https://tokei.rs/b1/github/grokify/entscape
 [repo-url]: https://github.com/grokify/entscape
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/entscape/blob/master/LICENSE
