# Design

## Context

See `proposal.md → Why` for motivation. Relevant existing patterns:

- **Artist** is the closest model: it has `last_used_at`, a unique-per-user name constraint, and a `UNIQUE(user_id, name)` index. Verse follows the same shape.
- **`artist.last_used_at` is bumped from the usecase layer** (not the repo), outside any transaction. The same best-effort approach applies to verse.
- **`characterUsecase`** currently depends on: `characterRepo`, `storage`, `avatarUploadRepo`, `transactor`, `tel`, `bookleafClient`. Adding `verseRepo` follows the same injection pattern.
- **`characterRepository.Create`** returns `error` and modifies the character struct in-place (no re-fetch at the end). This differs from `artistRepository.Create` which calls `GetByID` at the end. The existing pattern is kept — `verse_id` is a plain column, so it is already on the struct after assignment.

## Goals / Non-Goals

**Goals:**
- Introduce `verses` as a first-class entity with full CRUD and `last_used_at` tracking.
- Add `verse_id` as an optional FK on `characters`, assignable through existing create/update endpoints.
- Add `verse_id` query filter on `GET /characters`.
- Surface `verse_id` in character API responses.
- Keep all new code consistent with existing artist/character patterns — no new architectural layers.

**Non-Goals:**
- Frontend wiring (verse list UI, picker, character form integration) — deferred to a future pass.
- `GET /verses/:id` is included for completeness but has low FE priority right now.
- Many-to-many verse↔character — deliberately out of scope; one verse per character is the design choice.

## Decisions

### D1: Verse as a separate entity (not a tag/string)

Verses get their own table rather than being a freetext field on characters. This mirrors how artists are modeled and enables `last_used_at` recency tracking for a future picker. A plain string field could not support ordering by recent use.

Considered: freetext `verse` column on characters — rejected because it provides no recency signal and no CRUD surface.

### D2: `verse_id` FK with `ON DELETE SET NULL` (not CASCADE)

Consistent with every other "parent delete" decision in this schema: commission delete sets artpiece's `commission_id` to null, file delete does not cascade, artist delete sets artpiece/file `artist_id` to null. Verse delete follows the same pattern.

### D3: `last_used_at` bump from `characterUsecase`, outside the character transaction

The bump happens in the character usecase after a successful create/update, mirroring how `artpieceUsecase` and `commissionUsecase` bump `artist.last_used_at`. It is a best-effort recency signal; if the bump call fails after the character is persisted, the character is not rolled back. This is the accepted trade-off for keeping the operation simple.

Alternative considered: bump inside the character repository's transaction. Rejected because repositories in this codebase do not depend on each other; mixing verse and character operations in `characterRepository` would violate that boundary.

### D4: Verse ownership validated in `characterUsecase` before accepting `verse_id`

The FK `verse_id → verses(id)` enforces existence but not ownership (a user who knows another user's verse UUID could point to it). The usecase calls `verseRepo.GetByID(verseID, userID)` before proceeding — the same check `artpieceUsecase` applies for artist assignment. Returns `ErrVerseNotOwned` on failure, mapped to 422 in the handler.

### D5: `verse_id *uuid.UUID` on `domain.Character` (plain column, no association)

Since `characterRepository.Create` does not re-fetch, and the response only needs the UUID string (not a full verse embed), adding `Verse *domain.Verse` as a GORM association is unnecessary for this pass. A plain `VerseID *uuid.UUID` column on the domain struct is sufficient; the value flows through to the response as a UUID string. A full verse embed in the response is a future FE concern.

### D6: `verse_id` query param as `*string` with manual `uuid.Parse` in the handler

No existing query filter uses a UUID type. `ListCharacterFilters` uses `*string` for `q`; following the same approach for `verse_id` avoids introducing a new binding pattern. The handler parses and validates the UUID before passing `*uuid.UUID` to the usecase.

### D7: Two separate migrations

`000033_create_verses` creates the table. `000034_add_verse_id_to_characters` adds the FK column. Keeping them separate allows rolling back just the column addition without dropping the table, and matches the project's convention of one logical step per migration (see `000010` creating `artists`, `000024` adding `commission_id` to artpieces).

## Risks / Trade-offs

- **Best-effort `last_used_at` bump** — if the bump call fails after the character is persisted, the verse recency is stale. This is accepted for consistency with the artist pattern.
- **Verse_id filter on large character lists** — no index on `characters.verse_id` in the base migration. If filtering by verse becomes a hot path (unlikely at single-user scale), an index can be added in a follow-up migration.

## Migration Plan

1. Deploy migration `000033_create_verses` — new table, no impact on existing data.
2. Deploy migration `000034_add_verse_id_to_characters` — adds nullable column, all existing rows get `verse_id = NULL`. No backfill needed.
3. Deploy application code. The new `verse_id` field on character requests is optional; clients that don't send it continue to work without change.

Rollback: drop `verse_id` column (migration down), then drop `verses` table (migration down). No data loss on rollback since `verse_id` is nullable and characters remain intact.
