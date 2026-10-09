## Context

Artists currently support a single `artist_link` text field. The lookup dropdown (combobox used in commission and artpiece forms) sorts alphabetically. This change introduces a one-to-many `artist_links` table and a `last_used_at` recency column on artists, replacing both.

Existing `artist_link` values are intentionally not migrated — the field is dropped and users re-enter links going forward.

## Goals / Non-Goals

**Goals:**
- Allow an artist to have multiple links, one of which may be marked primary
- Surface the primary link wherever a single artist link was previously shown
- Sort the artist lookup dropdown by most-recently-used first

**Non-Goals:**
- Link labels (Discord, Twitter, VGen, etc.) — bare URLs only; can be added later as a nullable column with no migration pain
- Separate REST sub-resource endpoints for links (`POST /artists/:id/links`, etc.)
- Migrating existing `artist_link` values into the new table

## Decisions

### Links managed inline via replace-all, not sub-resource CRUD

Links are sent as a `links[]` array on the artist create/update body. On update, the backend deletes all existing links for that artist and inserts the new set.

**Why**: The existing commission-characters pattern already does this — `Association("Characters").Replace(...)` on every update. It keeps the API surface flat (no new routes), and the form modal has an explicit Save, so there is no need for live per-link mutations before commit.

**Alternative rejected**: Separate `POST/PATCH/DELETE /artists/:id/links` endpoints. Would require multiple round-trips on form submit and introduce a partial-save state (links saved, then artist name save fails). Inconsistent with every other form in the app.

### Links live in the existing `artistRepository` struct, not a separate injectable

Link DB operations (replace, preload) are implemented as methods on the existing `artistRepository`. No new interface is added to the usecase layer for links.

**Why**: Links are never accessed independently of their artist — they are always loaded as an association. Splitting them into a separate injectable would add wiring complexity for no architectural gain.

### Primary link enforced at the application layer, not a DB unique partial index

`is_primary = true` uniqueness per artist is checked/enforced by the application (only the first link with `is_primary: true` in the submitted array is honoured; extras are coerced to false). No `UNIQUE` partial index on `(artist_id) WHERE is_primary = true`.

**Why**: Consistent with how commission status and other single-active-value constraints are handled in this codebase. A partial unique index would make the replace-all pattern awkward (delete-then-insert ordering matters with constraints).

### `https://` prepend and URL validation on the frontend

On blur of each link input, the frontend normalises the value (`https://` prepended if no scheme) and validates with `new URL()`. Invalid links show an inline error below the input. The backend's `validate:"required,url"` on each link object remains as a safety net.

**Why**: Eliminates a round-trip for the most common mistake (typing `bsky.app/@artist` without a scheme). User sees the correction in place before submitting.

### `last_used_at` updated in the usecase layer, not via DB trigger

The artist usecase's `Create` sets `last_used_at = now()`. The commission and artpiece usecases call `artistRepo.UpdateLastUsedAt()` after validating artist ownership on create/update.

**Why**: Consistent with all other business logic in this codebase. DB triggers are opaque, bypass the usecase layer, and can't be unit-tested. The two narrow artist interfaces (`CommissionArtistRepository`, `ArtpieceArtistRepository`) each gain an `UpdateLastUsedAt` method — one concrete implementation on `artistRepository` satisfies both.

### Commission and artpiece API responses keep `artist_link: string | null`

The commission and artpiece response objects continue to expose `artist_link` as a flat `string | null`, resolved server-side to the primary link URL. The frontend consumers (`CommissionsTable`, `DashboardPage`) are unchanged.

**Why**: Avoids cascading a breaking change to commission/artpiece API consumers. The commission response already denormalises artist name as a flat field; primary link URL follows the same pattern.

## Risks / Trade-offs

**Existing link data is silently dropped** → Mitigation: documented in the migration plan; the links section is prominently visible on every artist edit modal, making re-entry obvious.

**All existing artists tie on `last_used_at` after backfill** → Mitigation: secondary sort key `name ASC` breaks ties alphabetically until real usage data accumulates. Acknowledged in the proposal.

**`last_used_at` is approximate under concurrent writes** → Two simultaneous commission creates assigning the same artist will both issue `UpdateLastUsedAt`; last write wins. Acceptable — sort order is a UX heuristic, not a precise record.

**`go-playground/validator` `dive` tag for slice validation** → Validates each `artistLinkInput` struct element including its `url` tag. The existing `validationErrResponse` reports field paths as `links[0].url`. No changes to the validator setup needed; `url` validation failures fall through to the generic "failed validation" message (current behaviour, not regressed).

## Migration Plan

1. Write migration `000029`:
   - Add `last_used_at timestamptz NULL` to `artists`
   - Create `artist_links` table: `id uuid PK`, `artist_id uuid NOT NULL REFERENCES artists ON DELETE CASCADE`, `url text NOT NULL`, `is_primary boolean NOT NULL DEFAULT false`, `created_at timestamptz`, `updated_at timestamptz`
   - Drop `artist_link` column from `artists`
   - Backfill: `UPDATE artists SET last_used_at = NOW()`
2. Deploy backend (new column nullable — no downtime concern)
3. Deploy frontend

Rollback: revert deploy; write a compensating migration that re-adds `artist_link` (nullable) and drops `artist_links`. Data loss on rollback is limited to any links users entered after the deploy — the original `artist_link` values were already dropped by the forward migration.
