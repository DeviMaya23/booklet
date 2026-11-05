# Design

## Context

See proposal.md for motivation. Current state, from the code:

- `CharacterFormModal` (531 lines) owns create and edit, the avatar pick/crop/upload logic, the folder diff wiring and an inline "Choose folders" column. `CharactersPage` always mounts it and toggles `open`, so `usePublicFolders()` fires when the Characters page loads.
- `ArtpieceDetailView` is the model for the new page: view/edit mode in one component, `unsavedChangesSummary()`, `useBlocker` + discard dialog, and a left-fields / right-tile layout. It is rendered by `ArtpiecesPage` when `useParams().id` is set.
- `GET /characters/:id` exists (404 on missing) but no FE hook reads it. `GET /artpieces` already returns `characters[]` and `is_reference`; the gallery filters it client-side.
- `useFolderDiff` takes `originalFolders` + the live public folders and tracks added/removed/renamed/invalid folders. Its reconciliation effect runs only when the live list changes.
- `components/ui/sheet.tsx` (Base UI dialog) exists; there is no tabs primitive. `AppSidebar` highlights with `useMatch(path, { end: false })`.
- `main.tsx` mounts `router.tsx`; `App.tsx` is not imported anywhere, so only the router needs the new route.

## Goals / Non-Goals

**Goals:**
- A character detail page that matches the mockup in view, edit and drawer states, with behavior parallel to the artpiece detail page.
- One implementation of avatar pick/crop/upload shared by the create modal and the page.
- Frontend only.

**Non-Goals:**
- Changing how characters are created (the modal keeps its fields and inline folder picker).
- Linking character chips elsewhere (e.g. in the artpiece detail) to the new page.
- Gallery view options (tile size, show details) on the tabs, pagination, or server-side filtering of artpieces by character.
- Showing Bookleaf folder images on the detail page itself; only the drawer previews them.

## Decisions

### D1. Separate page component and route
`/app/characters/:id` is added to `router.tsx` and rendered by a thin `pages/CharacterDetailPage.tsx` that reads the id and renders `features/characters/components/CharacterDetailView`. `CharactersPage` stays list-only.
*Alternative:* hang the id off `CharactersPage` like `ArtpiecesPage`. Rejected: it couples the list page's state (search, dialogs) to the detail view and the user chose a separate page. The URL shape is identical either way.

### D2. Split view and edit into two components
`CharacterDetailView` loads the character (new `useCharacter(id)`, key `['characters', 'detail', id]` so the existing `invalidateQueries(['characters'])` calls refresh it) and renders view mode. Clicking Edit mounts `CharacterEditForm`, which owns the draft state, the folder diff, the avatar picker and the unsaved-changes guard; Cancel or a successful Save unmounts it.
*Why not one component with a `mode` flag like `ArtpieceDetailView`:* `useFolderDiff` reconciles only when the live list changes. A long-lived component would need a manual `reset()` on Cancel/Save, after which reconciliation would not re-run and a no-longer-public folder could be silently re-included. Mounting per edit session makes reconciliation run fresh every time and removes the reset code.

### D3. Draft state and the unsaved-changes summary
`CharacterEditForm` holds name, notes and a verse override in local state (legitimate divergent state per CONVENTIONS). The verse is stored as `undefined` (untouched) / `Verse | null` (chosen), with the effective verse resolved from `useVerses()` against `character.verse_id`, so the form does not read as dirty while verses are still loading. A pure `characterEditSummary` function in `features/characters/lib` (colocated test; no React state, so not a hook) takes the draft, avatar picker state and selected folder IDs and returns the "Unsaved changes: …" string or `null`; `useBlocker` and the discard dialog key off that string, mirroring `ArtpieceDetailView`. Cancel discards directly; the back link and in-app navigation ask first.
Folder counts come from comparing selected IDs with the character's original IDs, so reconciliation-removed folders count as removed.

### D4. Shared avatar logic: one hook, one lib function
- `features/characters/hooks/useAvatarPicker` owns `localFile`, the preview object URL, `avatarCleared`, `cropSourceUrl`, drag-active state, size validation (toast), and the file/drop/paste/crop/clear handlers, with semantic mutators and URL revocation. The modal (create) and the page (edit) both use it.
- `features/characters/lib/uploadCharacterAvatar` runs init → `PUT` to the presigned URL → complete, parameterized by the init/complete mutate functions (the `uploadFile` pattern in `features/files/lib/upload.ts`). It replaces the modal's inline `runAvatarUpload`.
*Why now:* CONVENTIONS requires factoring a multi-step pipeline used from two entry points rather than copying it. The modal's create behavior is unchanged.
*Paste:* the modal keeps an `onPaste` on its dialog content. On the page no element has focus by default, so the page binds the hook's paste handler on `document` while `CharacterEditForm` is mounted; the existing "ignore paste from input/textarea" rule is inside the handler.

### D5. Save sequence and leaving the page
Save mirrors the modal's edit path: avatar upload or `DELETE …/avatar` first, then `PUT /characters/:id` (always sent, as today), then toast and unmount the form. Failures toast and keep the form mounted. After a successful delete or save the form must not trip its own blocker: the blocker predicate also checks a `leavingRef` set just before the post-delete `navigate('/app/characters')`.

### D6. Folders UI on the page
- View mode: `CharacterFolderChips` renders names, linking to `VITE_BOOKLEAF_HOST/share/:token` when `usePublicFolders()` has the token; on failure it renders plain names (no error UI).
- Edit mode: a two-column chip grid with count badge and description, warning rows for invalid folders, and a "Choose folders" button that opens `ChooseFoldersDrawer`, a thin `Sheet` wrapper around the existing `ChooseFoldersPanel`.
- The modal keeps `FolderList` untouched. The two chip lists differ in layout and affordances; per CONVENTIONS two feature-owned near-duplicates are acceptable until a third use or real drift appears.

### D7. Shared panel changes (modal and drawer)
`ChooseFoldersPanel` gains an `updatedAt` prop (from `dataUpdatedAt`) and shows "Updated <relative time>" using date-fns `formatDistanceToNow` (already a dependency), re-evaluated on a 60 s tick while open. `FolderPreview` changes `PAGE_SIZE` 15 → 24, adds the "Showing the first N" caption and a "View all in Bookleaf" link built from the folder token. These apply to the modal's inline panel too, which is intended: one panel, one behavior. The existing `FolderPreview`/`ChooseFoldersPanel` tests are updated for 24.

### D8. Tabs and data
- `components/ui/tabs.tsx` wraps `@base-ui/react` Tabs (`Tabs`, `TabsList`, `TabsTrigger`, `TabsContent`) in the style of `sheet.tsx`; triggers accept a count badge as children. No new dependency.
- `CharacterArtpieceTabs` takes the character id, reads `useArtpieces()`, filters by `characters.some(c => c.id === id)`, and partitions by `is_reference` in a `useMemo`. It reuses `ArtpiecesGrid` (fixed medium tiles, details shown) and `FileViewer`; `onInfoClick` navigates to `/app/artpieces/:id`.
- The `ArtpieceSummary → ViewerFile` mapping currently inlined in `ArtpiecesPage` is extracted to `features/artpieces/lib/toViewerFiles` and used by both pages (CONVENTIONS: avoid duplicated entry-point logic).
- `ArtpiecesGrid` renders "No artpieces yet." when empty, so the tab renders its own dashed empty-state block and only renders the grid when there are items.

### D9. CharactersPage and modal cleanup
`CharactersPage` drops `selectedCharacter`, the modal `key`, and `handleEditClick`; card click calls `navigate`. `CharacterFormModal` loses `character`, `isEditMode`, the Delete button/dialog, the avatar clear control, the edit submit path and the verse-sync effect, keeping the create flow and the discard-on-dirty dialog. The list page's card delete menu is unchanged.

## Risks / Trade-offs

- [Avatar committed but `PUT` fails] The avatar is already saved while the other fields are not → toast the failure, keep the form open, and invalidate the character queries after the avatar step so the page reflects the new avatar. Same ordering as the modal today.
- [Shared panel change alters the create modal] 24 previews in the narrower inline column → the grid is already 8 columns wide, so rows stay the same size; verify visually during apply.
- [Folder extra call in view mode] `GET /folders` now also fires on direct loads of the detail page → cached 15 minutes, `retry: false`, and view mode degrades to plain names.
- [Client-side partition of the full artpiece list] Cost grows with library size → acceptable for a single-user library and consistent with the gallery filter; revisit only if it shows.
- [Two similar folder chip lists] Drift risk between `FolderList` and `CharacterFolderChips` → accepted per CONVENTIONS; revisit if the modal later changes.
- [Unsaved-change guard vs. sidebar] `useBlocker` covers in-app navigation but not tab close or reload → same limitation as the artpiece detail.
