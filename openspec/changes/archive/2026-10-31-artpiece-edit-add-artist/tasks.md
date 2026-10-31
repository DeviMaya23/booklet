# Tasks

## 1. Add-artist action in edit mode

- [x] 1.1 In `ArtpieceDetailView.tsx`, add `artistModalOpen` state, wrap the edit-mode `ArtistCombobox` in a `flex gap-2` row (`flex-1` combobox + outline icon `Button` with `Plus`, `type="button"`, `disabled={isSaving}`), and render `<ArtistFormModal open onOpenChange onCreated={(artist) => setEditArtist(artist)} />`; verify by opening a detail page in edit mode, clicking "+", creating an artist, and seeing it selected in the combobox
- [x] 1.2 Verify the new button fits in the narrow artist column next to Characters at normal widths (visual check in the running app)

## 2. Tests

Follow the unit testing rules in `TESTING.md` / `CONVENTIONS.md`; mock `ArtistFormModal` the way `ArtpieceDetailView.test.tsx` already mocks its collaborators, if applicable.

- [x] 2.1 Test: clicking "+" in edit mode opens the artist-creation modal, and previously typed title/notes remain in the form after the modal is closed; verify the test passes
- [x] 2.2 Test: when the modal's `onCreated` fires with a new artist, the combobox shows that artist and the unsaved-changes summary contains "artist changed"; verify the test passes
- [x] 2.3 Test: after `onCreated`, clicking Save calls `mutateAsync` with `artistId` equal to the new artist's id; verify the test passes
- [x] 2.4 Test: closing the modal without creating leaves the previously selected artist unchanged and the unsaved-changes summary null; verify the test passes
- [x] 2.5 Test: the "+" button is disabled while saving; verify the test passes
- [x] 2.6 Run `npm test -- ArtpieceDetailView` in `frontend/` and verify all existing and new tests pass

## 3. Final checks

- [x] 3.1 Run `npm run build` in `frontend/` and fix any issue; verify it exits successfully
- [x] 3.2 Run `npm run lint` in `frontend/` and fix any issue; verify it exits with no errors
