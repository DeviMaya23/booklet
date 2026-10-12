# Tasks

## 1. Env var and type prep

- [x] 1.1 Add `VITE_BOOKLEAF_HOST=http://localhost:5173` to `frontend/.env.example`; verify it is present and `import.meta.env.VITE_BOOKLEAF_HOST` resolves at runtime
- [x] 1.2 Add `token: string` to the `PublicFolder` interface in `usePublicFolders.ts` and include it in the mapper; verify the field is present in query data

## 2. Modal header restructure

- [x] 2.1 Restructure `CharacterFormModal` to render the avatar thumbnail and Name field in a two-column row; verify the layout matches the design (small square thumbnail beside the Name input)
- [x] 2.2 Add a "Change" label overlay to the avatar thumbnail when an avatar exists; verify the overlay is visible in edit mode with an existing avatar and absent in create mode / when no avatar

## 3. Folder list (default view)

- [x] 3.1 Create `FolderList` component replacing `FolderPicker`/`TokenInput`; each valid folder row shows a folder icon, name, share link (`VITE_BOOKLEAF_HOST/share/:token`), and remove button; verify share link opens correct URL and remove triggers the existing `onRemove` callback
- [x] 3.2 Add `invalidFolders: PublicFolder[]` state to `CharacterFormModal`; update the reconciliation `useEffect` to populate `invalidFolders` alongside the existing `removedIds` logic (existing logic unchanged); verify both are reset in `resetForm`
- [x] 3.3 Render invalid folders as warning rows in `FolderList` (warning border, warning icon, "No longer public in Bookleaf.", Remove button); verify Remove dismisses the row (clears from `invalidFolders`) and the folder is already absent from the save payload
- [x] 3.4 Add "Choose folders" button below the folder list; verify clicking it sets `panelOpen` to true
- [x] 3.5 Write component tests for `FolderList`: render valid rows with share link, render warning row, remove valid folder, dismiss warning row

## 4. Choose folders panel

- [x] 4.1 Scaffold `ChooseFoldersPanel` as the right column of `CharacterFormModal`; Dialog widens to two-column layout when `panelOpen` is true; X button sets `panelOpen` to false; verify panel opens and closes correctly
- [x] 4.2 Render all public folders as checklist rows (checkbox, name, share link); checked state reflects current selection; checking adds the folder, unchecking removes it; verify both directions update the form's folder list
- [x] 4.3 Add client-side search bar; verify the list filters case-insensitively as the user types and clears when the input is cleared
- [x] 4.4 Add "Public folders from Bookleaf" subtitle and a refresh button that calls `publicFoldersQuery.refetch()`; verify refetch is triggered; verify a refetch failure retains the existing list (silent no-op)
- [x] 4.5 Add empty state when `publicFoldersQuery.data` is an empty array: "No public folders yet" message, "Open Bookleaf" link (opens `VITE_BOOKLEAF_HOST` in a new tab), Refresh button; verify empty state renders and Open Bookleaf has the correct href
- [x] 4.6 Add full-panel error state when `publicFoldersQuery.isError` and no data has loaded: "Couldn't reach Bookleaf" error block with Retry button; verify Retry calls `publicFoldersQuery.refetch()` and the folder list in the form remains functional during this state
- [x] 4.7 Write component tests for `ChooseFoldersPanel`: folder list renders with correct checked state, search filters, empty state, load-error state with Retry

## 5. Folder preview

- [x] 5.1 Create `FolderPreview` component; clicking Preview on a panel row sets `activePreviewId` to that folder's id and triggers `useGetFolderImages`; clicking Hide clears `activePreviewId`; verify only one preview is open at a time (opening a second closes the first)
- [x] 5.2 Show skeleton grid and "Loading images…" label while `useGetFolderImages` is fetching; verify skeleton renders during loading
- [x] 5.3 Render up to 15 thumbnails on load; show "Show more" button when more exist; each click appends 15 more from the full response; verify image count increments correctly
- [x] 5.4 Add per-folder error state when `useGetFolderImages` fails: inline "Couldn't load the images for this folder." message and Retry button; Retry calls the query's `refetch()`; verify error renders and Retry re-fetches
- [x] 5.5 Set `staleTime: 5 * 60 * 1000` on `useGetFolderImages` so repeated preview opens hit the cache; verify the query is not re-issued for a folder already fetched in the same session
- [x] 5.6 Write component tests for `FolderPreview`: skeleton during load, images rendered after load, show more increments count, error state with Retry

## 6. Build and lint

- [x] 6.1 Run `npm run build` from `frontend/`; fix any type or compilation errors
- [x] 6.2 Run `npm run lint` from `frontend/`; fix any lint errors without suppressing with `eslint-disable`

## Workflow follow-up

- Checkout a `feat/character-modal-revamp` branch from main before starting implementation.
- Archive the change once the implementation is reviewed and merged.
