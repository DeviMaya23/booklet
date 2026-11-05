# Tasks

## 1. Foundations

- [x] 1.1 Add `frontend/src/components/ui/tabs.tsx` wrapping `@base-ui/react` Tabs (`Tabs`, `TabsList`, `TabsTrigger`, `TabsContent`) in the style of `sheet.tsx`; verify `npm run build` type-checks it and the import resolves from `@base-ui/react/tabs`.
- [x] 1.2 Add `features/characters/api/useCharacter.ts` (`GET /characters/:id`, key `['characters', 'detail', id]`, disabled without an id) with a colocated test at the API boundary: success returns the character, a 404 puts the query in error, and `invalidateQueries(['characters'])` refetches it.
- [x] 1.3 Extract the `ArtpieceSummary → ViewerFile` mapping from `pages/ArtpiecesPage.tsx` into `features/artpieces/lib/toViewerFiles.ts`, use it in `ArtpiecesPage`, and add a unit test for the mapping (cover file id fallback, title fallback to cover name then "Untitled"); verify the existing `ArtpiecesPage.test.tsx` still passes.

## 2. Shared avatar logic

- [x] 2.1 Add `features/characters/lib/uploadCharacterAvatar.ts` (init → `PUT` presigned URL → complete, taking the init/complete mutate functions) with a test: the three steps run in order with the file's MIME type and size, and a non-OK `PUT` rejects before `complete` is called.
- [x] 2.2 Add `features/characters/hooks/useAvatarPicker.ts` (local file, preview URL, cleared flag, crop source URL, drag state, size-limit toast, file/drop/paste/crop/clear handlers, URL revocation, `reset`) with a colocated test: oversized file toasts and does not stage a crop, paste is ignored from inputs/textareas and without an image item, crop confirm sets the preview and clears the cleared flag, crop cancel changes nothing, and clear drops the local file.
- [x] 2.3 Switch `CharacterFormModal` to `useAvatarPicker` and `uploadCharacterAvatar`; verify the create flow is unchanged via a new `CharacterFormModal.test.tsx` covering the `web-characters` create scenarios (name only, with avatar, avatar failure deletes the new character and keeps the modal open, blank name disables Save).

## 3. Create-only modal and list navigation

- [x] 3.1 Remove the edit branch from `CharacterFormModal` (`character` prop, `isEditMode`, Delete button and dialog, avatar clear control, edit submit path, verse-sync effect); verify `npm run build` passes and the modal test from 2.3 asserts there is no Delete button and the title is "New character".
- [x] 3.2 Update `CharactersPage`: card click calls `navigate('/app/characters/:id')`, drop `selectedCharacter`/modal `key`/`handleEditClick`, keep the "New Character" button and card delete menu; add `CharactersPage.test.tsx` scenarios for card click navigating (no modal opens), the `(...)` menu not navigating, and "New Character" still opening the empty create modal.

## 4. Shared folder panel changes

- [x] 4.1 Update `FolderPreview` to 24 per page, add the "Showing the first N" caption and a "View all in Bookleaf" link (folder share URL, new tab); update `FolderPreview.test.tsx` for the 24/24 paging and add a test for the link `href`.
- [x] 4.2 Add an `updatedAt` prop to `ChooseFoldersPanel` showing "Updated <relative time>" (date-fns `formatDistanceToNow`, 60 s tick) and pass `dataUpdatedAt` from `usePublicFolders` in the modal; update `ChooseFoldersPanel.test.tsx` for the label (fake timers) and that Refresh calls `onRefresh`.

## 5. Detail page: route and view mode

- [x] 5.1 Add `pages/CharacterDetailPage.tsx` and the `/app/characters/:id` route in `router.tsx`; verify with a router test that the route renders the page and the Characters sidebar item is active on it.
- [x] 5.2 Build `CharacterDetailView` view mode: loading and "Failed to load character." states with back link, "← Characters" link, avatar/placeholder tile, name, verse chip (resolved from `useVerses`), notes, Edit button. Tests cover each `web-character-detail` view-mode scenario (header with data, no verse, no notes, no avatar, error state, back link).
- [x] 5.3 Add `CharacterFolderChips` view variant (chips, Bookleaf links when `usePublicFolders` has the token, plain names on failure, "No folders attached"); tests cover links present, `GET /folders` failing, and no folders.
- [x] 5.4 Add `CharacterArtpieceTabs` using `ui/tabs.tsx`, `useArtpieces`, `ArtpiecesGrid`, `FileViewer` and `toViewerFiles`: counts, the `is_reference` partition, exclusion of other characters' artpieces, dashed empty states with neutral copy, tile click opening the viewer, info control navigating to `/app/artpieces/:id`. Tests cover each tab scenario in the spec.

## 6. Detail page: edit mode

- [x] 6.1 Add `characterEditSummary` (pure function in `features/characters/lib`) with a colocated test: lists name/verse/notes/avatar changed, avatar removed, "N folder(s) added/removed" (including reconciliation-removed), returns `null` when nothing differs, and drops "verse changed" when the original verse is reselected.
- [x] 6.2 Build `CharacterEditForm` (mounted on Edit, unmounted on Cancel/Save): name input, Cancel/Save (Save disabled when the name is blank or saving), amber summary line, verse combobox with clear (verse override model), notes textarea, back link labelled with the name. Tests cover entering edit pre-filled, blank-name Save disabled, Cancel returning to view with original values, and the verse/notes request payloads (`verse_id: null`, `notes: null`).
- [x] 6.3 Wire avatar editing into `CharacterEditForm`: tile with Change overlay and X, click/drop, document-level paste while mounted, crop dialog, local-only preview and removal. Tests cover change via picker, oversized rejection, remove locally with no request before Save, and paste ignored in text fields.
- [x] 6.4 Add the edit-mode folder UI: `CharacterFolderChips` edit variant (count, description, ↗ link, remove, "No longer public in Bookleaf." warning rows with Remove), `ChooseFoldersDrawer` (Sheet around `ChooseFoldersPanel`, closes with X/Escape, edit state intact). Tests cover remove chip, warning row with exclusion from `folder_ids`, Bookleaf-unreachable chips, drawer open/close, drawer selections appearing as chips and "1 folder added".
- [x] 6.5 Implement save: avatar upload or `DELETE …/avatar`, then `PUT /characters/:id` with name, notes, `folder_ids`, `verse_id`; success toast and return to view with new values; failures toast and stay in edit; invalidate character queries after the avatar step. Tests cover metadata-only save, save with new avatar (order asserted), save after avatar removal (order asserted), avatar upload failure (no `PUT`), `PUT` failure (edit intact), and the diff `[A,B]` → `[A,C]`.
- [x] 6.6 Implement the unsaved-changes guard (`useBlocker` + "Discard changes?" dialog with a `leavingRef` bypass) and "Delete character" (confirm dialog, `DELETE`, toast, navigate to `/app/characters`, error toast on failure). Tests cover back-link and sidebar navigation guarded, Discard proceeding, Keep editing staying, no dialog without changes, delete confirmed/cancelled/failed, and no discard dialog after a successful delete with unsaved edits.

## 7. Review follow-ups

- [x] 7.1 Show loading and error states (not the empty state) in the artpiece tabs and hide the counts until data exists; verified by tab tests for a pending and a failed `GET /artpieces`.
- [x] 7.2 Write the PUT response into the detail query cache in `useUpdateCharacter` so view mode never flashes the old values; verified by `useUpdateCharacter.test.tsx` and the existing save test.
- [x] 7.3 Do not count Bookleaf-dropped folders as unsaved changes (summary ignores them, guard stays quiet) while still excluding them from the save; verified by summary unit tests and edit-form tests (no prompt on Back, user removal still counted, `folder_ids` unchanged).
- [x] 7.4 Split `CharacterEditForm` into `useUnsavedChangesGuard`, `useCharacterSave`, `useCharacterDelete`, a `listenForPaste` option on `useAvatarPicker`, and `CharacterAvatarEditor` / `DeleteCharacterDialog` / `DiscardChangesDialog`, each with a colocated test; verified by the unchanged edit-form behavior tests still passing.

## 8. Final checks

- [x] 8.1 Run `openspec validate character-detail-page --strict` and fix any findings.
- [x] 8.2 Run `npm run build` and `npm run lint` in `frontend/` and fix whatever they report; verify both exit 0 and the full Vitest suite passes.
- [x] 8.3 Run the app and compare view, edit and drawer states against the mockup (`character_edit_page_revamp.pdf`), plus the create modal's inline folder panel with the 24-image preview; verify no regressions in the create flow.
