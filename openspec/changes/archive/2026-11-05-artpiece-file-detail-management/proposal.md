# Proposal

## Why

Artpiece file management is currently image-focused: edit mode only exposes cover-setting and remove/add controls on thumbnail tiles, with no way to see or edit a file's `name`/`notes` from the artpiece context, and no way to keep an outdated file version attached without it cluttering the view. This change adds a dedicated file-details surface and a per-file visibility flag so outdated versions can stay attached without being shown by default.

## What Changes

- Add `is_visible` (bool, default `true`) to `File`, meaningful only when the file is attached to an artpiece (`artpiece_id` is set).
- Cover auto-selection (on attach, detach, create-with-files, and full-replace) skips files with `is_visible: false`.
- Explicitly hiding a file that is the artpiece's current cover triggers the same cover-reassignment logic as detaching it (first visible image, else first visible file, else null).
- `PUT /files/:id` accepts an optional `is_visible` field alongside the existing `name`/`notes` fields.
- The artpiece detail response's per-file entries gain `notes` and `is_visible` fields (previously only `name`, `width`, `height` were exposed there).
- New "Edit details" modal, opened from the artpiece detail view's files-section header (available in both view and edit mode): lists every attached file as a row (thumbnail, name input, notes input, visibility toggle), each field auto-saving independently on blur/toggle — no staging, no outer Save/Cancel involvement.
- Both the **view-mode** and **edit-mode** file grids exclude `is_visible: false` files entirely — hidden files never render as a tile in either grid. The files-section header shows a plain-text, non-interactive hint (e.g. "6 files · 2 hidden") when hidden files exist; there is no in-grid reveal control. The "Edit details" modal is the only surface where a hidden file's row appears.
- Edit-mode grid tiles otherwise keep their current controls — cover (☆) and remove (✕) stay staged/deferred-to-Save. Because hidden files never appear as tiles, there is no path to set a hidden file as cover from the grid.

## Capabilities

### New Capabilities
- `web-artpiece-file-details-modal`: the "Edit details" modal on the artpiece detail view — lists attached files with auto-saving name, notes, and visibility controls per row.

### Modified Capabilities
- `artpiece-files`: cover auto-selection (attach/detach/create/full-replace) now skips hidden files; hiding the current cover reassigns it via the existing detach-reassignment rules; the artpiece file response gains `notes` and `is_visible`.
- `file-crud`: `PUT /files/:id` accepts an optional `is_visible` field; file responses include `is_visible`.
- `web-artpiece-detail`: view-mode file grid filters out hidden files by default with a toggle to reveal them; files-section header gains an "Edit details" entry point in both view and edit mode.

## Impact

- **Backend**: `domain.File` gains `IsVisible bool`; a migration adds the `is_visible` column (`not null default true`); `file_usecase.Update` signature changes to accept `is_visible`; `artpiece_usecase` cover-selection logic (used by attach, detach, create-with-file_ids, and full-replace-files) filters by visibility; a new usecase path handles "hide the current cover" reassignment; `toFileResponse` and the artpiece-detail file-serialization path both gain fields.
- **Frontend**: new `ArtpieceFileDetailsModal` component; `ArtpieceFile` type gains `notes`/`is_visible`; `ArtpieceDetailFileGrid` (view mode) filters by visibility with a reveal toggle; `ArtpieceDetailView` gains the modal's open state and entry point, reachable from both view and edit mode.
- **Bruno**: new/updated request files for `PUT /files/:id` (is_visible field).
