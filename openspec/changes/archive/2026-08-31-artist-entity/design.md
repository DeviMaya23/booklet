## Context

Images currently embed artist attribution as two free-text columns (`artist_name`, `artist_link`) on both `images` and `pending_uploads`. This requires re-entry on every upload and makes artist information non-reusable. The change promotes artist into a first-class entity owned per-user, referenced by FK from images and pending uploads.

## Goals / Non-Goals

**Goals:**
- Introduce an `artists` table with CRUD endpoints scoped to the owning user
- Replace inline artist text on `images` and `pending_uploads` with a nullable `artist_id` FK
- Ensure artist deletion cascades gracefully (SET NULL, not block or cascade delete)
- Keep the upload flow accepting `artist_id` at initiation; validate existence at completion

**Non-Goals:**
- Shared/global artist registry across users
- Migrating existing `artist_name`/`artist_link` data — dropped by migration
- Image creation flow changes beyond accepting `artist_id` (FE flow is a separate proposal)
- Soft delete on artists

## Decisions

### Artist entity shape
`artists(id UUID PK, user_id UUID NOT NULL FK, name TEXT NOT NULL, notes TEXT, artist_link TEXT, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ)`

Uniqueness: `UNIQUE(user_id, name)` composite constraint — same name is allowed for different users, not for the same user.

No `deleted_at` / soft delete. Artist is hard-deleted; the FK constraint handles referencing rows via `ON DELETE SET NULL`.

`artist_link` (not `url`) — consistent with the existing `artist_link` naming on `images`.

### FK behavior on delete
Both `images.artist_id` and `pending_uploads.artist_id` use `ON DELETE SET NULL`. Rationale: artist is informational metadata. Deleting an artist shouldn't block deletion or cascade-destroy images.

### CompleteUpload artist validation
At `CompleteUpload` time, if `pending.ArtistID` is non-nil, the usecase queries for the artist by `(id, user_id)`. If not found (artist was deleted between initiation and completion), `artist_id` on the created image is silently set to `nil`. This mirrors the existing pattern for character IDs (silently filtered if invalid).

**Alternative considered**: rely purely on `ON DELETE SET NULL` at DB level and skip the usecase check. Rejected because the DB constraint only fires on artist deletion, not on `pending_uploads.artist_id` being stale — the pending row could hold an ID from a different user's artist if validation were skipped at initiation. Explicit usecase validation is clearer.

### No `artist_id` validation at InitialUpload
`artist_id` is accepted as an optional UUID at initiation (shape-validated only, like `character_ids`). Existence is not checked until `CompleteUpload`. Rationale: consistent with how character IDs are handled; avoids a race condition where an artist deleted between validation and persistence would cause a spurious 422.

### PATCH /images/:id artist field
`updateImageRequest` replaces `artist_name *string` + `artist_link *string` with `artist_id Patch[string]`. PATCH semantics: if `artist_id` is absent from the body, the existing association is unchanged. If present (including `null`), it replaces the current value.

`Patch[string]` is a generic wrapper that distinguishes three states: absent (`.Set == false`), explicit null (`.Set == true`, `.Value == nil`), and a concrete value (`.Set == true`, `.Value != nil`). This is necessary because Go's `*string` collapses absent and null into the same `nil` value, making proper PATCH semantics impossible.

**Alternative considered**: `json.RawMessage` for deferred decoding. Rejected in favour of `Patch[string]` — the generic wrapper is type-safe, reusable, and integrates cleanly with the go-playground validator via a registered custom type func.

### Artist ownership validation on image update
When `PATCH /images/:id` sets a non-null `artist_id`, the image repository (or usecase) must verify the artist belongs to the requesting user before applying the update, returning 422 if not. Prevents a user from associating another user's artist with their image.

## Risks / Trade-offs

- **Breaking API change** → Clients reading `artist_name`/`artist_link` from image responses will break. No mitigation — accepted as a known breaking change with no live clients beyond internal tools at this stage.
- **Data loss on migration** → Existing `artist_name`/`artist_link` text is dropped. No backfill. Accepted; data volume is low and the fields were free-text with no relational value.
- **Artist deleted mid-upload window** → Handled by silent null at CompleteUpload. The user won't see an error but the image will land with no artist. Acceptable given the 15-minute presign TTL and the informational nature of artist data.

## Migration Plan

Single migration (`000010_add_artists_entity`):
1. Create `artists` table
2. `ALTER TABLE images DROP COLUMN artist_name, DROP COLUMN artist_link, ADD COLUMN artist_id UUID REFERENCES artists(id) ON DELETE SET NULL`
3. `ALTER TABLE pending_uploads DROP COLUMN artist_name, DROP COLUMN artist_link, ADD COLUMN artist_id UUID REFERENCES artists(id) ON DELETE SET NULL`

Rollback: down migration drops `artists`, removes `artist_id` columns, re-adds `artist_name`/`artist_link` (empty — data is unrecoverable).

No deploy coordination needed beyond deploying the migration before the new binary; old binary writes `artist_name`/`artist_link` which the migration drops, so deploy order is: migrate first, then deploy binary.
