## Context

Files currently have a `notes` field as the only user-visible text. The `UpdateNotes` method is a narrow per-field updater at every layer (repo interface, usecase, handler interface). Adding `name` as a second mutable text field creates a choice: add a parallel `UpdateName` method, or generalise the existing one.

The upload pipeline has two stages: `InitiateUpload` creates a `pending_file_uploads` row, then `CompleteUpload` creates the `files` row from it. Any field that must survive that handoff must live in `pending_file_uploads` as well.

## Goals / Non-Goals

**Goals:**
- Add `name TEXT` (nullable) to `files` and `pending_file_uploads`
- Make `name` settable at upload initiation and updatable via PATCH
- Replace the two-call pattern that would result from keeping `UpdateNotes` separate with a single atomic `Update(name, notes)`
- Surface `name` in the API response and extend inbox search to include it

**Non-Goals:**
- Enforcing uniqueness or any validation on `name`
- Back-filling `name` on existing file rows
- Exposing `name` as a server-side search/filter parameter (client-side only)

## Decisions

### Decision: Rename `UpdateNotes` → `Update(name, notes *string)` at all layers

**Chosen**: Expand the single update method to accept both `name` and `notes` as nullable pointers. The repository issues one `UPDATE files SET name = ?, notes = ? WHERE id = ? AND user_id = ?`. The usecase and handler interface follow the same change.

**Alternative considered**: Add a separate `UpdateName` method and call both from the handler. Rejected because it produces two sequential DB writes with no atomicity guarantee — if `UpdateNotes` succeeds and `UpdateName` fails, the row is left in a partially-updated state. It also doubles the surface area (two methods, two test paths) for what is conceptually one operation.

### Decision: `name` is passed at `InitiateUpload`, not `CompleteUpload`

**Chosen**: The frontend has access to the browser `File` object (including `.name`) only in the drag-and-drop handler, which calls `InitiateUpload`. `CompleteUpload` has no request body. Storing `name` in `pending_file_uploads` and copying it to `files` on completion is the natural flow.

**Alternative considered**: Adding a body to `CompleteUpload`. Rejected — unnecessary protocol change when the data is already available at initiation time.

### Decision: `name` is nullable with no default

**Chosen**: No back-fill, no server-side default. Existing files get `NULL` name; new uploads via DnD set it from the browser filename. This keeps the migration trivial (`ALTER TABLE ... ADD COLUMN name TEXT`) and avoids touching existing data.

## Risks / Trade-offs

- **`UpdateNotes` rename touches many files** (handler, usecase, repo, and their tests) → The change is mechanical; a grep for `UpdateNotes` before starting implementation confirms all sites, which the explore phase already mapped.
- **`name` set from browser `File.name` can be anything** (long path, unicode, `.DS_Store`, etc.) → Acceptable for now; there is no validation requirement and the field is nullable/display-only.
- **`pending_file_uploads` rows older than the presign TTL are purged** regardless of `name` — no impact since the purge path does not use `name`.

## Migration Plan

1. Deploy migration `000021` that adds `name TEXT` to both tables — additive, no downtime
2. Deploy backend with updated domain structs, usecase, handler
3. Deploy frontend with updated `useInitFileUpload` signature and search filter
4. No rollback complexity: `name` is nullable, old code ignores the column, new code handles NULL gracefully
