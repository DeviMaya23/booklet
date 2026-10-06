## 1. Shared Utilities

- [x] 1.1 Create `frontend/src/lib/files.ts` with `stripFileExtension(name: string): string` — strips the last extension from a filename, handles dotfiles and trailing dots per the design rules
- [x] 1.2 Add `validateFileName(name: string): string | null` to `frontend/src/lib/files.ts` — returns an error string for blank, forbidden chars (`/ \ : * ? " < > |`), or > 255 chars; null if valid

## 2. Upload Call Sites — Extension Strip

- [x] 2.1 `AddFilesModal.tsx`: apply `stripFileExtension` to `file.name` before passing to `initUpload.mutateAsync({ name })` and before setting the uploaded row's initial `name`
- [x] 2.2 `ArtpieceFilesInput.tsx`: same — apply `stripFileExtension` at both the `initUpload` call and the uploaded-state `name` assignment
- [x] 2.3 `FileInboxGrid.tsx`: apply `stripFileExtension` to `file.name` before passing to `initUpload.mutateAsync({ name })`
- [x] 2.4 `ArtpieceDetailView.tsx`: apply `stripFileExtension` to `file.name` before passing to `initUpload.mutateAsync({ name })`

## 3. Name Input Validation

- [x] 3.1 `AddFilesModal.tsx`: call `validateFileName` in `handleFieldBlur` for the name field; if it returns an error, set `nameError` to that message and skip the `updateFile` call
- [x] 3.2 `ArtpieceFilesInput.tsx`: add `nameError` state to the uploaded `DraftFile` variant; call `validateFileName` in `handleNameBlur`; show inline error below the name input (red border + error text, matching `AddFilesModal` style)

## 4. Linting

- [x] 4.1 Run `npm run build` and `npm run lint` in `frontend/`; fix any issues
