## Context

Booklet's backend already has a user provisioning flow (middleware → usecase → repository) but no application-level CRUD has been implemented yet. `characters` is the first entity with full handler → usecase → repository wiring. The patterns established here will serve as the reference for future entities.

## Goals / Non-Goals

**Goals:**
- Introduce the `Character` domain entity and its five CRUD operations
- Enforce ownership at the repository layer (all queries filter by `user_id`)
- Establish the handler → usecase → repository pattern that future entities will follow
- Handle partial updates (PATCH) correctly for all field types, including booleans

**Non-Goals:**
- Public/unauthenticated read access (`is_public` is stored but not acted upon here)
- Pagination or filtering on the list endpoint
- Image upload to R2 (hero_image_r2_path is plain text for now)

## Decisions

### UUID generation at the usecase layer

Character IDs are UUIDs (`github.com/google/uuid`, already in go.mod). Generation happens in the usecase before the record is handed to the repository — not at the database level and not in the handler.

_Rationale_: Keeps ID generation testable and co-located with business logic. The repository stays a pure persistence layer.

### Ownership enforced at the repository layer

Every query (GetByID, Update, Delete) includes `WHERE id = ? AND user_id = ?`. A wrong-owner lookup returns `gorm.ErrRecordNotFound`, which the handler maps to 404 — indistinguishable from a missing record. This avoids leaking ownership information.

_Alternative considered_: Enforce in usecase by fetching first then comparing. Rejected — requires an extra round-trip and the DB is the authoritative source of the row.

### PATCH with pointer fields + map conversion

`UpdateCharacterParams` uses `*string` / `*bool` fields. Nil means "field was omitted from the request body." The repository converts non-nil pointers to a `map[string]interface{}` before passing to GORM's `.Updates()`.

_Rationale_: Pointer fields give nil/non-nil as the "was-this-sent?" signal, which handles the `is_public: false` case correctly. GORM's `.Updates(struct)` skips zero values (false, "", 0), so the map conversion is necessary to prevent silently ignoring an explicit `false`.

### Repository interface in its own file

`CharacterRepository` has 5 methods, which exceeds the ≤2-method threshold for inlining. It lives in `usecase/character_repository.go` per conventions.

### Soft delete via GORM DeletedAt

`Character` embeds `gorm.DeletedAt`. Deleted records are automatically excluded from all standard GORM queries. No explicit `WHERE deleted_at IS NULL` clauses needed.

### Error wrapping with `%w`

Repository methods wrap errors with `fmt.Errorf("...: %w", err)`, preserving `gorm.ErrRecordNotFound` for `errors.Is` checks in the handler layer. This matches the existing `user_repository.go` pattern.

## Risks / Trade-offs

- **No pagination on list**: The list endpoint returns all of a user's characters. Acceptable for MVP; could become a problem at scale. → Add pagination when needed; the usecase signature can evolve.
- **hero_image_r2_path is a raw string**: No validation that the path exists in R2. → Path integrity will be enforced when the upload flow is implemented.
- **Soft delete only**: There is no hard-delete or restore endpoint. Soft-deleted records accumulate. → Add a purge job when needed.

## Migration Plan

1. Create migration `000002_create_characters.up.sql` — adds `characters` table with FK to `users`
2. Create migration `000002_create_characters.down.sql` — drops the table
3. Run `make migrate-up` in the target environment before deploying the new binary
4. Rollback: run `make migrate-down` then redeploy previous binary
