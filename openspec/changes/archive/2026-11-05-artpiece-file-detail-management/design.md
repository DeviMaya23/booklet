# Design

## Context

See `proposal.md` — Why. Relevant current state:

- `domain.File` has no visibility concept; `artpiece_id` nullable already distinguishes "in the Inbox" vs "attached."
- Cover auto-selection lives entirely in `ArtpieceUsecase` (`backend/internal/usecase/artpiece_usecase.go`): `maybeAutoSetCover` (on attach), `reassignCover` (on detach), and `setCoverFromFileSet` (on create-with-`file_ids` and full-replace). All three currently pick "first image file, else first file" with no visibility filter.
- `PUT /files/:id` (`FileHandler.UpdateFile` → `FileUsecase.Update`) already takes `name *string, notes *string`.
- The artpiece-detail serialization path (handler-level, separate from `fileResponse` used by `file-crud` endpoints) currently emits `name`, `extension`, `width`, `height` per attached file — not `notes`.
- `ArtpieceDetailFileGrid.tsx` renders one shared grid for both view and edit mode, driven by a discriminated-union `mode` prop; `FileTileEditOverlay.tsx` and `AddFilesModal.tsx` are the two existing precedents for an auto-save, per-row file-metadata UI.

## Goals / Non-Goals

**Goals:**
- Make hiding a file's visibility fully consistent, regardless of entry point (explicit toggle, or implicitly by never letting a hidden file become cover).
- Keep the new modal's save model (immediate, per-field) distinct from the artpiece edit session's save model (staged, batch-committed) without entangling the two.

**Non-Goals:**
- No change to what counts as an "image" file for cover purposes (`bookmime.IsImage`) — visibility is an orthogonal filter on top of the existing image-preference rule.
- No backend restriction on `SetCover` (`PUT /artpieces/:id/cover`) explicitly targeting a hidden file. Because hidden files never render as tiles in either grid, the UI has no path to call it with a hidden file's ID; adding a server-side guard for a path the UI can't reach would be speculative.
- No change to `is_visible` when a file is detached — the value is simply inert while `artpiece_id` is null, and re-attaching a previously-hidden file honestly preserves the state the user last set.

## Decisions

**`is_visible bool not null default true` on `files`, via a numbered migration.** Matches the existing nullable-`artpiece_id` pattern of "a plain column that's only semantically active in one context." No `NULL` state needed since the attribute is inert when the file isn't attached — `NOT NULL DEFAULT true` keeps every existing row trivially valid.

**Visibility filtering is centralized at the three existing cover-selection call sites, not duplicated per endpoint.** `maybeAutoSetCover`, `reassignCover`, and `setCoverFromFileSet` all change to filter their candidate list to `IsVisible` files before applying the existing image-preference logic. This is a small, local change to logic that's already the single source of truth for cover selection — no new usecase method needed for attach/detach/create/full-replace.

**Hiding the current cover reuses `reassignCover`, not a new algorithm.** `FileUsecase.Update`, when it detects `is_visible` flipping to `false` on a file that is some artpiece's current `cover_file_id`, calls into the same cover-reassignment path `DetachFile` already uses. The candidate list for that call is naturally already visibility-filtered under the decision above, so the newly-hidden file excludes itself without any extra exclusion parameter.

**`FileUsecase.Update` needs the file's artpiece association to decide whether reassignment applies**, so this logic sits in `FileUsecase` but calls out to the artpiece-side reassignment function (an existing cross-usecase call, same shape as other usecases already composing on each other in this codebase) — not duplicated as file-side logic.

**The "Edit details" modal is a new, separate component, not a mode on `ArtpieceDetailFileGrid`.** It has its own save model (immediate, per-field, matching `FileTileEditOverlay`/`AddFilesModal`) that's fundamentally different from the grid's staged/batched model. Reusing the grid component for both would require threading two incompatible save strategies through one component's props.

**The modal fetches the artpiece's file list from the already-loaded `useArtpiece` query data** (now including `notes` and `is_visible` per the `artpiece-files` delta) rather than issuing a separate fetch — it's opened from a view that already has this data in cache, and auto-save mutations invalidate that same query key so the grid and modal stay in sync.

**Both grids (view and edit mode) filter by `is_visible` client-side**, the same way `ArtpieceDetailFileGrid` already branches on `mode` — no new backend endpoint for "visible-only" is needed since the full file list (needed by the modal) is already in the response.

## Risks / Trade-offs

**[Risk]** Toggling visibility in the modal changes what the edit-mode grid shows (a tile can appear/disappear) while the user has an *unrelated* staged edit session open (e.g. mid-editing the title). → **Mitigation**: this is already true today for uploads ("Uploaded files persist through Cancel" — immediate, not staged), so it's consistent with an existing precedent rather than a new class of behavior. No change needed to the Cancel/unsaved-changes logic.

**[Risk]** A file that's hidden while it's the only remaining visible candidate during a *bulk* operation (e.g. full-replace swaps in several files, all hidden) could leave an artpiece coverless where a naive implementation might expect a cover. → **Mitigation**: explicitly specified as a scenario in the `artpiece-files` delta ("New set contains only hidden files" / "Create with all specified files hidden") — coverless is the correct, intended outcome, not a bug.

**[Risk]** Migration adds a `NOT NULL` column with a default — safe for existing rows, but worth confirming the migration tool backfills via the column default rather than requiring an explicit `UPDATE`. → **Mitigation**: standard `ADD COLUMN ... NOT NULL DEFAULT true` is backfilled by Postgres without a table rewrite lock in modern versions; call this out explicitly in the migration's `up.sql` task so it isn't accidentally written as a two-step add-then-backfill.

## Migration Plan

Standard additive migration: new nullable-free column with a default, no data backfill script needed, no breaking change to any existing endpoint (new fields are additive on responses; `is_visible` on `PUT /files/:id` is optional). Rollback is a straightforward column drop in `down.sql`.
