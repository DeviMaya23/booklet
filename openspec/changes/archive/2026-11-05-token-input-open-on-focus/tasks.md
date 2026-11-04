# Tasks

## 1. Shared chip input cap

- [x] 1.1 In `TokenInput.tsx`, slice the empty-text list to 10 unselected suggestions and render a non-selectable hint row when truncated; keep it out of the keyboard row count. Verify: new `TokenInput.test.tsx` cases pass (25 items shows 10 + hint; 10 or fewer shows no hint; selected items don't consume the cap; typing finds an item beyond the first 10; arrow keys and clicking the hint add nothing).

## 2. Characters on focus

- [x] 2.1 Add `openOnFocus` to the Characters `TokenInput` in `ArtpieceFormModal.tsx`. Verify: modal test shows recency-ordered characters on focus and the capped list with hint; confirm the Inbox entry (`FilesPage`) and `ArtpiecesPage` both render this modal.
- [x] 2.2 Add `openOnFocus` to the Characters `TokenInput` in `ArtpieceDetailView.tsx` (edit mode). Verify: detail view test shows characters on focus, excludes selected chips, and shows the capped list with hint.
- [x] 2.3 Add a `TagPicker` test for the capped focus list. Verify: `TagPicker.test.tsx` passes.

## 3. Final checks

- [x] 3.1 Run `npm run build` and `npm run lint` in `frontend/` and fix any issues. Verify: both exit cleanly.
