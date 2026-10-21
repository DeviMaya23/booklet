# Proposal

## Why

Saved filters exist and can be created from the artpiece gallery, but the dashboard shows no trace of them. This is Part 2: wire saved filters into the dashboard with display, navigation, per-filter actions, and visibility control.

## What Changes

- **`ThumbnailTile` component** (new, shared): a navigation-first square thumbnail card with label, optional sublabel, and an optional `menuItems` action menu. Used by both dashboard sections. Replaces the inline `RecentArtpiecesTile` in `DashboardPage` and serves as the tile for the saved filters section.
- **Recent artpieces tiles**: wired up click-to-navigate to artpiece detail (`/app/artpieces/:id`) via `ThumbnailTile`.
- **`ResourceCard`**: unchanged. Remains characters-only (edit-modal primary action, delete-only menu).
- **Saved filters section on dashboard**: horizontal tile row above Recent artpieces. Each tile shows the saved filter's thumbnail and name. Clicking navigates to `/app/artpieces` with the filter pre-applied via URL params. A `...` menu per tile exposes "Change thumbnail" and "Delete".
- **"Choose shown" toggle**: a "Choose shown" button on the section header opens a popover listing all saved filters with checkboxes. Visibility state is stored in `localStorage` (`dashboard.hiddenFilterIds`). All filters are shown by default; hiding is opt-out. The section shows an empty state when all filters are hidden.
- **Delete saved filter**: `...` menu item triggers a confirmation `AlertDialog`. On confirm, calls `DELETE /saved_filters/:id`. On success, the tile is removed from the row.
- **Change thumbnail**: `...` menu item opens a modal. The modal loads `useArtpieces` (cached), applies the saved filter's `filter_payload` defensively, and renders a picker grid (up to 10 artpieces). First artpiece is pre-selected. On save, calls `PATCH /saved_filters/:id` with `thumbnail_artpiece_id`. On success, the tile thumbnail updates.
- **Backend PATCH update**: `PATCH /saved_filters/:id` gains `thumbnail_artpiece_id` support alongside the existing `thumbnail_r2_path` field. The handler resolves the artpiece's cover thumbnail R2 path, same as the CREATE handler.
- **New FE API hooks**: `useSavedFilters` (list query), `useDeleteSavedFilter` (mutation), `useUpdateSavedFilter` (mutation).

## Capabilities

### New Capabilities

- `web-thumbnail-tile`: the shared `ThumbnailTile` component — navigation-first square thumbnail card with optional action menu, used across dashboard sections.
- `web-dashboard-saved-filters`: the saved filters section on the dashboard — tile row, "Choose shown" visibility toggle, per-tile delete and change-thumbnail actions, and click-to-navigate behavior.

### Modified Capabilities

- `web-dashboard`: adds the saved filters section requirement above recent artpieces; adds click-to-navigate behavior to recent artpieces tiles.
- `saved-filters`: `PATCH /saved_filters/:id` gains `thumbnail_artpiece_id` as an accepted field for updating the thumbnail.

## Impact

- **Backend**: `saved_filter_handler.go` — `patchSavedFilterRequest` gains `ThumbnailArtpieceID *string`; patch handler gains artpiece resolution logic (reuse from create handler).
- **Frontend**: new `components/ThumbnailTile.tsx`.
- **Frontend**: `DashboardPage.tsx` — imports `ThumbnailTile`; removes inline `RecentArtpiecesTile`; adds saved filters section with `useSavedFilters`, choose-shown popover, `ChangeThumbnailModal`, and `DeleteSavedFilterDialog`.
- **Frontend**: new `features/artpieces/components/ChangeThumbnailModal.tsx`.
- **Frontend**: new `features/artpieces/api/useSavedFilters.ts`, `useDeleteSavedFilter.ts`, `useUpdateSavedFilter.ts`.
- **Bruno**: update `collection/saved_filters/` patch file to show `thumbnail_artpiece_id` in request.
- No route changes. No new pages.
