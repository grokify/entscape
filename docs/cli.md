# CLI Reference

## Commands

### entscape html

Generate a standalone HTML visualization file.

```bash
entscape html [schema-dir] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `schema-dir` | Path to Ent schema directory (e.g., `./ent/schema`) |

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file path (default: stdout) |
| `--repo` | | Repository URL for source links |
| `--branch` | | Git branch for source links (default: `main`) |
| `--title` | | Page title (default: package name or "Entscape") |

**Examples:**

```bash
# Basic usage
entscape html ./ent/schema -o index.html

# With source links
entscape html ./ent/schema \
  -o docs/erd.html \
  --repo https://github.com/org/repo \
  --branch main \
  --title "My App"
```

---

### entscape generate

Parse Ent schema and export as JSON.

```bash
entscape generate [schema-dir] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `schema-dir` | Path to Ent schema directory |

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file path (default: stdout) |
| `--repo` | | Repository URL for source links |
| `--branch` | | Git branch for source links (default: `main`) |
| `--docs` | | Documentation base URL |
| `--pretty` | | Pretty-print JSON (default: `true`) |

**Examples:**

```bash
# Export to stdout
entscape generate ./ent/schema

# Export to file
entscape generate ./ent/schema -o schema.json

# With source links
entscape generate ./ent/schema \
  -o schema.json \
  --repo https://github.com/org/repo
```

---

### entscape serve

Start a local development server.

```bash
entscape serve [schema-dir] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `schema-dir` | Path to Ent schema directory |

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--port` | `-p` | Port to serve on (default: `8080`) |
| `--repo` | | Repository URL for source links |
| `--branch` | | Git branch for source links (default: `main`) |
| `--docs` | | Documentation base URL |
| `--web` | | Web directory to serve static files from |

**Endpoints:**

| Endpoint | Description |
|----------|-------------|
| `/api/schema.json` | Parsed schema as JSON |
| `/api/jsonschema` | JSON Schema definition |
| `/health` | Health check |

**Examples:**

```bash
# Basic server
entscape serve ./ent/schema

# Custom port
entscape serve ./ent/schema -p 3000

# With web UI
entscape serve ./ent/schema --web web
```

---

### entscape schema

Generate JSON Schema for the entscape format.

```bash
entscape schema [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file path (default: stdout) |

**Examples:**

```bash
# Output to stdout
entscape schema

# Save to file
entscape schema -o entscape.schema.json
```

---

### entscape version

Print version information.

```bash
entscape version
```

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | Error (parsing, file I/O, etc.) |

## Environment Variables

Currently, entscape does not use environment variables. All configuration is done via command-line flags.
