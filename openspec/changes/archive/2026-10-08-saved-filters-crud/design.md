## Context

The project has a well-established pattern for user-owned CRUD resources (Artist, Character, Commission): domain struct → repository interface → usecase → handler → route registration. This change adds `SavedFilter` as a new resource in that same mold, with one new wrinkle: the first `jsonb` column in any domain struct.

## Goals / Non-Goals

**Goals:**
- Full CRUD for `saved_filters` following existing layered architecture patterns
- Opaque storage of `filter_payload` as a jsonb blob; no backend validation of its internal shape

**Non-Goals:**
- Frontend integration (separate change)
- Pagination on the list endpoint — return all records for the authenticated user
- Unique constraint on `name` — duplicate names are allowed
- Soft delete — hard deletion only

## Decisions

### D1: `json.RawMessage` for `filter_payload`

**Decision**: Use `json.RawMessage` for the `filter_payload` domain field and response struct field.

**Rationale**: GORM handles `json.RawMessage` transparently for `jsonb` columns (read/write without extra marshaling). No new dependency needed (`gorm.io/datatypes` would be the alternative but adds a dep for one field). `json.RawMessage` also embeds inline in JSON responses — no double-encoding.

**Alternative considered**: `[]byte` — same effect but less idiomatic; `json.RawMessage` is the standard stdlib type for this purpose.

### D2: `*json.RawMessage` (not `Patch[json.RawMessage]`) in PATCH request

**Decision**: Use `*json.RawMessage` for `filter_payload` in the PATCH request struct.

**Rationale**: `filter_payload` is `NOT NULL` in the DB — it can never be explicitly cleared to null. Three-state semantics (`Patch[T]`) exist to distinguish "absent" from "explicit null clear"; that distinction is meaningless for a NOT NULL field. `nil` pointer = field absent (skip update), non-nil = replace payload. The existing validator only registers `Patch[string]`; using `*json.RawMessage` avoids needing to extend that registration.

### D3: `Patch[string]` for `name` and `thumbnail_r2_path` in PATCH

**Decision**: Use `Patch[string]` for `name` (to validate non-empty when provided) and `thumbnail_r2_path` (nullable — needs three-state to allow explicit clearing).

**Rationale**: Follows the established `Patch[T]` pattern already used in the codebase. The validator's custom type func already handles `Patch[string]`.

### D4: Hard delete, no soft delete

**Decision**: No `gorm.DeletedAt` on `SavedFilter`.

**Rationale**: Saved filters are user-facing bookmarks with no recovery scenario. Artist and Commission also use hard delete. Soft delete is only on Character; that's specific to that entity's needs.

### D5: List ordering

**Decision**: `ORDER BY created_at DESC` on `GET /saved_filters`.

**Rationale**: Most recently created filters are most likely to be the ones the user wants. No client-specified ordering.

## Risks / Trade-offs

- **Opaque payload drift**: Old saved filters may reference artist/character IDs that no longer exist. The backend stores and returns the blob faithfully; the frontend is responsible for handling stale references at read time. This is an explicit design decision, not an oversight.
- **No payload size cap**: A malicious or buggy client could send a very large `filter_payload`. Low risk given this is an authenticated endpoint for personal use, but worth noting.

## Migration Plan

- Migration `000028_create_saved_filters` — single `CREATE TABLE` with no data backfill needed.
- Rollback: `DROP TABLE saved_filters` — no FK dependencies from other tables.
