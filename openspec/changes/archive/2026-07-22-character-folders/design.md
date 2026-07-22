## Context

Characters are the core entity in Booklet. Bookleaf is a separate service that manages folder collections. Users organise content into Bookleaf folders, and characters need to be associated with those folders. The folder data (name, token, ID) is owned by Bookleaf; Booklet only stores the folder IDs that are assigned to each character.

The current character module (`domain`, `usecase`, `repository`, `handler`) supports CRUD with no folder concept. This change adds the association layer and a proxy endpoint for folder discovery.

## Goals / Non-Goals

**Goals:**
- Store folder–character associations in a local `character_folders` join table
- Return assigned folder IDs on all character read responses
- Accept optional folder ID assignments on character create and update
- Provide a `GET /folders` endpoint that proxies the Bookleaf public-folders API for the authenticated user
- Validate that provided folder IDs are UUID-shaped (no existence check against Bookleaf)

**Non-Goals:**
- Validating that a folder ID exists in Bookleaf before assigning it (deferred)
- Storing or using the `token` or `folder_name` fields from the Bookleaf response beyond proxying them
- Pagination on the folders endpoint

## Decisions

### 1. `character_folders` table with composite PK, no surrogate key

The join table uses `(character_id, folder_id)` as the composite primary key. No surrogate `id` column, no `created_at`. A character cannot be assigned the same folder twice; uniqueness is enforced by the PK.

`folder_id` is `uuid` type but has no FK constraint — it references an external system. `character_id` is a FK to `characters.id`.

**Alternative considered**: storing `folder_id` as `text` to avoid implying a DB relationship. Rejected — uuid column type gives free format validation at the DB layer and is consistent with the rest of the schema.

### 2. Repo-level transaction orchestration for folder writes

Folders are an attribute of character, not a separate aggregate. The character repository owns both `characters` and `character_folders` and manages them atomically using an internal DB transaction. The usecase passes `FolderIDs *[]uuid.UUID` in params and is unaware of the join table mechanics.

This avoids introducing `Transactor` into the character usecase, which is warranted only when multiple repositories are involved (as in the upload flow).

**Alternative considered**: usecase-level orchestration with `Transactor`. Rejected — folders are a character attribute, and leaking the join table boundary into the usecase adds complexity without benefit.

### 3. Update semantics: nil = no-op, empty slice = clear all

`FolderIDs *[]uuid.UUID` in both `CreateCharacterParams` and `UpdateCharacterParams`:
- `nil` → no change to folder assignments (or none on create)
- `&[]uuid.UUID{}` → remove all folder assignments
- `&[]uuid.UUID{id1, id2}` → replace current assignments with this set

On update, the repo performs a replace-all: delete existing `character_folders` rows for the character, then insert the new set, in a single transaction.

### 4. Hard-delete folder rows on character soft-delete

Characters are soft-deleted via GORM's `DeletedAt`. `character_folders` rows have no soft-delete semantics — they are hard-deleted when their character is deleted (in the same transaction). The folder rows are invisible once the character is soft-deleted, but accumulating orphans is unnecessary.

### 5. Bookleaf client in `internal/bookleaf/`

The Bookleaf HTTP client is a thin adapter over an internal REST API. It does not belong in `internal/storage/` (which is scoped to object storage concerns). A dedicated `internal/bookleaf/` package makes the dependency explicit and namespaced.

The client is constructed with the base URL from `BOOKLEAF_HOST` and the internal secret from `BOOKLEAF_INTERNAL_SECRET` (both required env vars). Every request to the Bookleaf API includes the header `X-Bookleaf-Internal-Secret: <secret>`. The usecase defines the interface; `internal/bookleaf/` implements it.

**`BOOKLEAF_HOST` and `BOOKLEAF_INTERNAL_SECRET` are both required env vars** — missing either prevents server startup, consistent with the pattern for all required config in this codebase.

### 6. `GET /folders` is a pure proxy — no usecase transformation

The folder listing endpoint calls the Bookleaf API and returns the response body as-is. The response shape (`folder_list` with `folder_id`, `token`, `folder_name`) is passed through without transformation. The handler calls the folder usecase; the usecase calls the bookleaf client and returns its result unchanged.

Per conventions, no unit test is written for the usecase happy path (pure delegation). The handler test covers the happy path (HTTP plumbing) and the error mapping (500 on upstream failure).

## Risks / Trade-offs

- **Bookleaf availability** → `GET /folders` fails if Bookleaf is down. Mitigation: standard 500 response; no retry or caching in scope. Users see an error and can retry.
- **Stale folder assignments** → A folder deleted from Bookleaf is not removed from `character_folders`. Mitigation: deferred (folder validity check is out of scope). Orphaned assignments are benign until that endpoint is built.
- **Replace-all on update is destructive** → A client that sends a partial folder list will lose unincluded assignments. Mitigation: documented in API spec; the client is responsible for sending the complete desired set.
- **No cascade from `characters` FK** → The FK from `character_folders.character_id` to `characters.id` does not use `ON DELETE CASCADE` (characters are soft-deleted, not hard-deleted). Mitigation: the repository hard-deletes folder rows explicitly in the delete transaction.
