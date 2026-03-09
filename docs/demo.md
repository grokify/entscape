# Demo

This page shows a live demo of the entscape visualization using a sample Ent schema.

## Interactive ERD

<div style="width: 100%; height: 600px; border: 1px solid #dee2e6; border-radius: 8px; overflow: hidden;">
  <iframe src="erd/index.html" style="width: 100%; height: 100%; border: none;"></iframe>
</div>

!!! tip "Interaction Tips"
    - **Click** an entity to open its source file
    - **Hover** over an entity to see its fields
    - Use the **toolbar** to change layout direction
    - **Scroll** to zoom in/out

## Sample Schema

The demo uses a sample schema with 5 entities demonstrating various relationship types:

| Entity | Description | Relationships |
|--------|-------------|---------------|
| **User** | User accounts | Has many Posts, has one Profile, belongs to many Groups |
| **Post** | Blog posts | Belongs to User, has many Comments |
| **Profile** | User profiles | Belongs to User (1:1) |
| **Comment** | Post comments | Belongs to Post |
| **Group** | User groups | Has many Users (M:M) |

## Source Code

The sample schema files are available in the repository:

- [user.go](https://github.com/grokify/entscape/blob/main/testdata/basic/schema/user.go)
- [post.go](https://github.com/grokify/entscape/blob/main/testdata/basic/schema/post.go)
- [profile.go](https://github.com/grokify/entscape/blob/main/testdata/basic/schema/profile.go)
- [comment.go](https://github.com/grokify/entscape/blob/main/testdata/basic/schema/comment.go)
- [group.go](https://github.com/grokify/entscape/blob/main/testdata/basic/schema/group.go)

## Generate Your Own

```bash
entscape html ./ent/schema \
  -o docs/erd/index.html \
  --repo https://github.com/your-org/your-repo \
  --title "Your Schema"
```
