## 1. AddFilesModal component

- [x] 1.1 Create `frontend/src/features/files/components/AddFilesModal.tsx` with a `Dialog` shell, title "Add files to dump", and a Close button footer
- [x] 1.2 Add the drop zone button (drag-and-drop + file picker trigger) at the top of the modal body, matching the visual style of `ArtpieceFilesInput`'s drop zone
- [x] 1.3 Add local `DraftFile` state union type: `uploading` (clientId, fileName) | `uploaded` (clientId, id, name, notes, thumbnailUrl, thumbnailGenState, nameError, notesError)
- [x] 1.4 Implement `uploadFiles` using `useInitFileUpload` and `useCompleteFileUpload`: add `uploading` rows immediately, transition to `uploaded` on success, remove row + toast on failure
- [x] 1.5 Wire the drop zone's `onDrop` and file input's `onChange` to `uploadFiles`
- [x] 1.6 Sync `thumbnailUrl` and `thumbnailGenState` into uploaded rows via a `useEffect` on the `useFiles` query result (same pattern as `ArtpieceFilesInput`)
- [x] 1.7 Render uploading rows: spinner thumbnail, disabled name input (filename), disabled notes input, no X button
- [x] 1.8 Render uploaded rows: thumbnail slot (image / spinner / file icon based on `thumbnailGenState`), enabled name input, enabled notes input, X button; disable name and notes while `thumbnailGenState === 'pending'`
- [x] 1.9 Implement name/notes blur handlers: call `useUpdateFile`, clear the relevant error on success, set `nameError`/`notesError` on failure
- [x] 1.10 Render inline field errors: red border on the affected input + a short error message below it when `nameError` or `notesError` is set
- [x] 1.11 Implement X button: remove the row from local state (file stays in dump)

## 2. Wire into FilesPage

- [x] 2.1 Add `addFilesModalOpen` state to `FilesPage` and pass it to `AddFilesModal`
- [x] 2.2 Wire the floating "+" button `onClick` to open the modal

## 3. Quality

- [x] 3.1 Run `npm run build` from the frontend directory and fix any type errors
- [x] 3.2 Run `npm run lint` and fix any lint issues
