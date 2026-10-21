# Tasks

## 1. Backend: PATCH thumbnail_artpiece_id

- [x] 1.1 Add `ThumbnailArtpieceID *string` to `patchSavedFilterRequest`; in the handler, when set, fetch the artpiece via `ArtpieceRepository.GetByID` and resolve its cover thumbnail R2 path (same logic as `CreateSavedFilter`), populate `UpdateSavedFilterParams.ThumbnailR2Path`; when both `thumbnail_artpiece_id` and `thumbnail_r2_path` are present, `thumbnail_artpiece_id` takes precedence; verify handler compiles
- [x] 1.2 Add handler unit tests: PATCH with `thumbnail_artpiece_id` for an artpiece that has a cover thumbnail returns 200 with non-null `thumbnail_url`; PATCH with `thumbnail_artpiece_id` for an artpiece with no cover thumbnail returns 200 with null `thumbnail_url`; verify tests pass
- [x] 1.3 Update `collection/saved_filters/patch_saved_filter.bru` to include `thumbnail_artpiece_id` in the request body example; verify file reflects the updated shape

## 2. Frontend: API hooks

- [x] 2.1 Create `features/artpieces/api/useSavedFilters.ts` — GET `/saved_filters` list query returning `SavedFilter[]` (reuse the existing `SavedFilter` type from `useCreateSavedFilter`); verify it compiles with no TypeScript errors
- [x] 2.2 Create `features/artpieces/api/useDeleteSavedFilter.ts` — DELETE `/saved_filters/:id` mutation; verify it compiles
- [x] 2.3 Create `features/artpieces/api/useUpdateSavedFilter.ts` — PATCH `/saved_filters/:id` mutation accepting `{ id: string; thumbnailArtpieceId: string }` and sending `thumbnail_artpiece_id` in the body; returns updated `SavedFilter`; verify it compiles

## 3. Frontend: parseFilterPayload utility and ThumbnailTile component

- [x] 3.1 Create `features/artpieces/lib/parseFilterPayload.ts` — defensive parser for raw `filter_payload` JSON; returns `{ artist_id: string | null; character_ids: string[]; match_mode: 'all' | 'any' }` or `null` on any shape mismatch; verify file exists and compiles
- [x] 3.2 Write unit tests for `parseFilterPayload`: valid full payload, null `artist_id`, empty `character_ids`, non-object input, missing fields, `match_mode` defaults to `'all'` when absent or unrecognised; verify all pass
- [x] 3.3 Create `components/ThumbnailTile.tsx` — `<Link>`-wrapped square image area with muted placeholder on no URL or load error, label below, optional sublabel, optional `menuItems` dropdown (absent → no `...` button); verify it renders and TypeScript compiles
- [x] 3.4 Write unit tests for `ThumbnailTile`: image renders when URL provided; placeholder renders when no URL; no action button when `menuItems` absent; action button opens menu when `menuItems` provided; destructive item has destructive style; menu item click does not trigger navigation; verify all pass

## 4. Frontend: Dashboard — tile migration and saved filters section

- [x] 4.1 Remove inline `RecentArtpiecesTile` from `DashboardPage`; replace with `ThumbnailTile` using `href={/app/artpieces/${a.id}}`; verify the dashboard renders recent artpiece tiles and clicking one navigates to `/app/artpieces/:id`
- [x] 4.2 Add `useSavedFilters` to `DashboardPage`; render the Saved Filters section above Recent artpieces using one `ThumbnailTile` per filter; clicking a tile navigates to `/app/artpieces` with URL params derived from `parseFilterPayload` (`artist`, `characters`, `match`); a `null` parse result navigates to `/app/artpieces` with no params; verify section renders and navigation produces correct URL params
- [x] 4.3 Wire the Delete action: `...` menu "Delete" item opens an `AlertDialog` confirmation; confirming calls `useDeleteSavedFilter` and invalidates the `useSavedFilters` query on success; cancelling sends no request; verify tile is removed on confirm and unchanged on cancel
- [x] 4.4 Add "Choose shown" button to the Saved Filters section header; it opens a popover listing all saved filters with a checkbox each; toggling persists hidden IDs to `localStorage` key `dashboard.hiddenFilterIds` (try/catch on all access); hidden filters do not render as tiles; when all filters are hidden the section shows an empty state message; verify toggle state persists across page reload

## 5. Frontend: ChangeThumbnailModal

- [x] 5.1 Create `features/artpieces/components/ChangeThumbnailModal.tsx` — shows the filter summary chip (built from parsed `filter_payload`), a picker grid of up to 10 artpieces from cached `useArtpieces` filtered by `filterAndSortArtpieces`, first artpiece pre-selected, placeholder on image load error, empty state with disabled Save when no artpieces match; on Save calls `useUpdateSavedFilter` with selected artpiece ID; verify component compiles
- [x] 5.2 Wire `ChangeThumbnailModal` into the saved filter tile `...` menu "Change thumbnail" item in `DashboardPage`; on save success invalidate the `useSavedFilters` query so the tile thumbnail updates; verify modal opens from the menu, saving patches the filter and the tile reflects the new thumbnail
- [x] 5.3 Write unit tests for `ChangeThumbnailModal`: grid renders up to 10 filtered artpieces; first artpiece is pre-selected; empty state renders and Save is disabled when no artpieces match the filter; clicking Cancel closes without calling `useUpdateSavedFilter`; verify all pass

## 6. Quality

- [x] 6.1 Run `npm run build` and fix any type or build errors
- [x] 6.2 Run `npm run lint` and fix any lint errors
- [x] 6.3 Run `golangci-lint run ./...` and fix any issues

## Workflow follow-up

- Archive the change once implementation is complete and reviewed.
