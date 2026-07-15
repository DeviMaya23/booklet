## Why

Booklet's core purpose is showcasing original characters — but there is currently no way to create or manage them. This change introduces the `characters` entity with full CRUD, establishing the first handler → usecase → repository flow in the codebase.

## What Changes

- New `characters` table in PostgreSQL (migration)
- New domain entity `Character` with soft-delete support
- Full CRUD API under `/characters` (protected, owner-scoped):
  - `POST /characters` — create a character
  - `GET /characters/:id` — get a single character by ID
  - `GET /characters` — list all characters owned by the authenticated user
  - `PATCH /characters/:id` — partial update (name, hero_image_r2_path, biography, is_public)
  - `DELETE /characters/:id` — soft delete
- Character IDs are UUIDs generated at the usecase layer
- All access is scoped to the authenticated user; ownership enforced at the repository layer via `user_id` filter
- Bruno collection files for all five endpoints

## Capabilities

### New Capabilities

- `character-management`: Create, read, list, update (partial), and soft-delete characters owned by the authenticated user

### Modified Capabilities

_(none)_

## Impact

- **New files**: `domain/character.go`, `usecase/character_usecase.go`, `usecase/character_repository.go`, `repository/character_repository.go`, `handler/character_handler.go`, migration files, Bruno collection files
- **Modified files**: `cmd/server/main.go` — wire up new repository, usecase, handler, and register routes on the `protected` group
- **Dependencies**: `github.com/google/uuid` (already in go.mod)
- **No breaking changes** to existing endpoints or contracts
