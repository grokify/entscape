# Getting Started

## Installation

### Go Install

```bash
go install github.com/grokify/entscape/cmd/entscape@latest
```

### From Source

```bash
git clone https://github.com/grokify/entscape.git
cd entscape
go build -o entscape ./cmd/entscape
```

## Basic Usage

### Generate HTML Visualization

The simplest way to visualize your Ent schema:

```bash
entscape html ./ent/schema -o index.html
```

This creates a standalone HTML file that you can open directly in a browser.

### Add Source Links

Link entities to your source code:

```bash
entscape html ./ent/schema \
  -o index.html \
  --repo https://github.com/your-org/your-repo \
  --branch main \
  --title "My App Schema"
```

### Export JSON

Export the parsed schema as JSON for custom tooling:

```bash
entscape generate ./ent/schema -o schema.json --pretty
```

### Development Server

Start a local server for development:

```bash
entscape serve ./ent/schema --port 8080
```

The server provides:

- `/api/schema.json` - Parsed schema as JSON
- `/api/jsonschema` - JSON Schema definition
- `/health` - Health check endpoint

## GitHub Pages Deployment

### Option 1: Manual

```bash
# Generate HTML
entscape html ./ent/schema \
  -o docs/index.html \
  --repo https://github.com/your-org/your-repo

# Commit and push
git add docs/
git commit -m "docs: add ERD visualization"
git push

# Enable GitHub Pages in repo settings (source: docs/ folder)
```

### Option 2: GitHub Actions

Create `.github/workflows/erd.yml`:

```yaml
name: Generate ERD

on:
  push:
    paths:
      - 'ent/schema/**'
    branches:
      - main

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install entscape
        run: go install github.com/grokify/entscape/cmd/entscape@latest

      - name: Generate ERD
        run: |
          entscape html ./ent/schema \
            -o docs/index.html \
            --repo ${{ github.server_url }}/${{ github.repository }} \
            --title "${{ github.repository }}"

      - name: Commit changes
        run: |
          git config user.name github-actions
          git config user.email github-actions@github.com
          git add docs/
          git diff --staged --quiet || git commit -m "docs: update ERD"
          git push
```

## NPM Module

For custom integrations, use the NPM module:

```bash
npm install @grokify/entscape
```

```javascript
import { render } from '@grokify/entscape';

// Fetch schema from your API
const response = await fetch('/api/schema.json');
const schema = await response.json();

// Render visualization
const instance = render('#container', schema, {
  direction: 'TB',
  onEntityClick: (entity) => {
    console.log('Clicked:', entity.name);
  }
});

// Instance methods
instance.fit();
instance.highlight('User');
instance.zoomTo('Post');
```
