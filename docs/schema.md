# JSON Schema

Entscape exports schema data in a standardized JSON format. This page documents the schema structure.

## Schema Version

Current version: `1`

## Root Object

```json
{
  "version": "1",
  "package": { ... },
  "entities": [ ... ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Schema version |
| `package` | object | Package metadata (optional) |
| `entities` | array | List of entities |

## Package

```json
{
  "name": "github.com/org/repo/ent",
  "source": "https://github.com/org/repo",
  "branch": "main",
  "docs": "https://pkg.go.dev/github.com/org/repo/ent"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Go package path |
| `source` | string | Repository URL |
| `branch` | string | Git branch |
| `docs` | string | Documentation URL |

## Entity

```json
{
  "name": "User",
  "path": "https://github.com/org/repo/blob/main/ent/schema/user.go",
  "description": "User entity",
  "fields": [ ... ],
  "edges": [ ... ],
  "indexes": [ ... ],
  "mixins": [ ... ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Entity name |
| `path` | string | Source file URL |
| `description` | string | Entity description |
| `fields` | array | Field definitions |
| `edges` | array | Edge (relationship) definitions |
| `indexes` | array | Index definitions |
| `mixins` | array | Mixin names |

## Field

```json
{
  "name": "email",
  "type": "string",
  "attrs": ["required", "unique"],
  "default": "",
  "comment": "User's email address"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Field name |
| `type` | string | Go type (string, int, time.Time, etc.) |
| `attrs` | array | Field attributes |
| `default` | string | Default value |
| `comment` | string | Field comment |

### Field Attributes

| Attribute | Description |
|-----------|-------------|
| `primary` | Primary key |
| `unique` | Unique constraint |
| `required` | Not nullable |
| `optional` | Nullable |
| `immutable` | Cannot be updated |
| `sensitive` | Sensitive data (excluded from logs) |
| `nillable` | Pointer type in Go |
| `default` | Has default value |

## Edge

```json
{
  "name": "posts",
  "target": "Post",
  "relation": "O2M",
  "inverse": "author",
  "required": false,
  "unique": false,
  "comment": "User's posts"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Edge name |
| `target` | string | Target entity name |
| `relation` | string | Relation type |
| `inverse` | string | Inverse edge name |
| `required` | boolean | Required relationship |
| `unique` | boolean | Unique relationship |
| `comment` | string | Edge comment |

### Relation Types

| Type | Description |
|------|-------------|
| `O2O` | One-to-One |
| `O2M` | One-to-Many |
| `M2O` | Many-to-One |
| `M2M` | Many-to-Many |

## Index

```json
{
  "fields": ["email"],
  "unique": true
}
```

| Field | Type | Description |
|-------|------|-------------|
| `fields` | array | Field names in the index |
| `unique` | boolean | Unique index |

## Full Example

```json
{
  "version": "1",
  "package": {
    "name": "github.com/example/app/ent",
    "source": "https://github.com/example/app",
    "branch": "main"
  },
  "entities": [
    {
      "name": "User",
      "path": "https://github.com/example/app/blob/main/ent/schema/user.go",
      "fields": [
        {
          "name": "id",
          "type": "int",
          "attrs": ["primary"]
        },
        {
          "name": "email",
          "type": "string",
          "attrs": ["required", "unique"]
        },
        {
          "name": "name",
          "type": "string",
          "attrs": ["required"]
        },
        {
          "name": "created_at",
          "type": "time.Time",
          "attrs": ["immutable", "default"]
        }
      ],
      "edges": [
        {
          "name": "posts",
          "target": "Post",
          "relation": "O2M",
          "inverse": "author"
        }
      ],
      "indexes": [
        {
          "fields": ["email"],
          "unique": true
        }
      ]
    },
    {
      "name": "Post",
      "path": "https://github.com/example/app/blob/main/ent/schema/post.go",
      "fields": [
        {
          "name": "id",
          "type": "int",
          "attrs": ["primary"]
        },
        {
          "name": "title",
          "type": "string",
          "attrs": ["required"]
        },
        {
          "name": "content",
          "type": "string"
        }
      ],
      "edges": [
        {
          "name": "author",
          "target": "User",
          "relation": "M2O",
          "inverse": "posts",
          "required": true
        }
      ]
    }
  ]
}
```
