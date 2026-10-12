# Proposal

## Why

The character modal's folder picker uses a chip/token input that gives no feedback about invalid folders, no way to preview folder contents before linking, and no direct path to Bookleaf. The overall modal layout also buries the avatar behind a full-width hero zone that doesn't use space well.

## What Changes

- **Modal layout:** Avatar thumbnail moves inline with the Name field (two-column row). Avatar area shows a "Change" overlay label; clear button remains.
- **Folder list (default/collapsed state):** Replaces the `TokenInput` chip picker with a plain list of rows. Each valid row shows folder name, a Bookleaf share link (↗), and a remove button. Folders silently auto-removed during reconciliation are now surfaced as explicit warning rows ("No longer public in Bookleaf.") with a manual Remove button — the underlying silent-removal behaviour is unchanged; this is a visibility layer only.
- **"Choose folders" panel (expanded state):** A side panel slides in beside the form. It lists all public Bookleaf folders with checkboxes (checked = currently assigned). Each row has a Preview button (on-demand image grid, one open at a time), a share link, and a toggle to hide the preview. Supports client-side search across the folder list.
- **Folder preview:** Calls existing `GET /folders/:id/images` endpoint. API returns all images; FE shows 15 initially, +15 per "Show more" (client-side). No "Showing the first N" label. No "View all in Bookleaf" link (the per-row ↗ covers it).
- **Empty state:** When the folder list is empty, the panel shows "No public folders yet" with an "Open Bookleaf" CTA and a Refresh button.
- **Error states:**
  - Initial folder list load failure: full-panel error block ("Couldn't reach Bookleaf"), Retry.
  - Preview image load failure: inline per-folder error + Retry.
  - Folder list refetch failure while data already loaded: silent no-op.
- **Bookleaf share links:** `VITE_BOOKLEAF_HOST/share/:token`. Token is already returned in `GET /folders` but currently dropped by the FE mapper — add `token` to `PublicFolder` type and mapper.
- **New env var:** `VITE_BOOKLEAF_HOST` added to `.env.example` (value: `http://localhost:5173`).
- **Timestamp removed:** "Updated X ago" label omitted; only "Public folders from Bookleaf" + refresh icon shown.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `folder-linking-ui`: folder picker UI is being fully replaced; panel, preview, error states, warning rows, and share links are new behaviour; degraded-state and reconciliation requirements change
- `web-characters`: `CharacterFormModal` layout changes (avatar + name row); folder field presentation changes

## Impact

- `frontend/src/features/characters/components/CharacterFormModal.tsx` — layout, folder state, panel open/close
- `frontend/src/features/characters/components/FolderPicker.tsx` — replaced with new list + panel components
- `frontend/src/features/characters/api/usePublicFolders.ts` — add `token` to `PublicFolder` type and mapper
- `frontend/src/features/characters/api/useGetFolderImages.ts` — consumed by preview; no changes expected
- `frontend/.env.example` — add `VITE_BOOKLEAF_HOST`
- No backend changes required
