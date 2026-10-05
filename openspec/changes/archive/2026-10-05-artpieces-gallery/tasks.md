## 1. ResourceCard sublabel

- [x] 1.1 Add optional `sublabel?: string` prop to `ResourceCard` — when present, render it as a smaller line below the existing label in the gradient overlay; existing callers are unaffected

## 2. ArtpiecesGrid component

- [x] 2.1 Create `features/artpieces/components/ArtpiecesGrid.tsx` — card grid using `ResourceCard`, passing `label` as `title ?? "Untitled"` and `sublabel` as `artist_name` (omit if null)

## 3. Filter popover component

- [x] 3.1 Create `features/artpieces/components/ArtpiecesFilterPopover.tsx` with artist single-select `Combobox` (suggestions from `useArtists()`), character `TokenInput` (suggestions from `useCharacters()`), and All/Any segmented toggle defaulting to "All"
- [x] 3.2 Implement outside-click dismissal via `useEffect` document `mousedown` listener
- [x] 3.3 Derive and expose active filter summary string (e.g. `"1 artist"`, `"2 characters, all"`, `"1 artist · 2 characters, any"`) for the Filter button label

## 4. ArtpiecesPage

- [x] 4.1 Create `pages/ArtpiecesPage.tsx` with toolbar row (search input, Filter button+popover, sort dropdown) and `ArtpiecesGrid`; pin "+ New Artpiece" button to the right of the toolbar
- [x] 4.2 Wire client-side search: filter by `title` case-insensitive substring; artpieces with null title are excluded when search is non-empty
- [x] 4.3 Wire client-side sort: "Newest first" (default, by `created_at` desc) and "Alphabetical" (by `title` asc, null titles last)
- [x] 4.4 Wire client-side artist filter: show only artpieces whose `artist_id` matches the selected artist
- [x] 4.5 Wire client-side character filter: "All" mode requires every selected character ID to appear in the artpiece's `characters` array; "Any" mode requires at least one
- [x] 4.6 Wire "+ New Artpiece" button to open `ArtpieceFormModal` with no `initialFileIds`; on `onSuccess` invalidate `ARTPIECES_QUERY_KEY`

## 5. Routing and navigation

- [x] 5.1 Add `/app/artpieces` route in `App.tsx` rendering `ArtpiecesPage`
- [x] 5.2 Add "Artpieces" entry to `navItems` in `AppSidebar` between Inbox and Characters, linking to `/app/artpieces`

## 6. Delete artpiece

- [x] 6.1 Update `web-artpieces-gallery` spec to add delete requirement and scenarios
- [x] 6.2 Create `features/artpieces/api/useDeleteArtpiece.ts` — mutation calling `DELETE /artpieces/:id`, invalidates `ARTPIECES_QUERY_KEY` on success
- [x] 6.3 Wire delete into `ArtpiecesGrid`: confirmation `AlertDialog`, success toast "Artpiece deleted", error toast "Failed to delete artpiece"

## 7. QA

- [x] 7.1 Run `npm run build` in `frontend/`, fix any type or build errors
- [x] 7.2 Run `npm run lint` in `frontend/`, fix any lint errors
