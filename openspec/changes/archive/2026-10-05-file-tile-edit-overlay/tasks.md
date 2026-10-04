## 1. FileTile updates

- [x] 1.1 Add a truncated read-only filename label beneath the tile image; hide it when `file.name` is null
- [x] 1.2 Remove the `···` `DropdownMenu` and "View Detail" item from `FileTile`
- [x] 1.3 Add a hover trash icon button to the tile that calls `onDeleteClick`; style consistently with the existing hover chrome
- [x] 1.4 Wire `onDoubleClick` prop on the tile's root div (already passed in from `FileInboxGrid` as a no-op; remove the no-op and let the parent supply the real handler)

## 2. FileTileEditOverlay component

- [x] 2.1 Create `frontend/src/features/files/components/FileTileEditOverlay.tsx` — fixed-position overlay panel covering the right portion of the grid, with a backdrop that closes the overlay on click; closes on Escape key
- [x] 2.2 Render the file's thumbnail (or fallback icon) at the top of the overlay
- [x] 2.3 Add name and notes `Input` fields pre-filled from the file prop
- [x] 2.4 Implement `handleFieldBlur(field: 'name' | 'notes')` — calls `useUpdateFile`, clears the field's error on success, sets `"Couldn't save. Try again."` on failure
- [x] 2.5 Render inline field errors: `border-destructive` on the input + short error message below when error is set

## 3. FileInboxGrid wiring

- [x] 3.1 Add `editingFile` state (`File | null`) to `FileInboxGrid`
- [x] 3.2 Pass `onDoubleClick={() => setEditingFile(file)}` to each `FileTile`
- [x] 3.3 Mount `<FileTileEditOverlay>` when `editingFile` is non-null, passing the file and `onClose={() => setEditingFile(null)}`

## 4. Quality

- [x] 4.1 Run `npm run build` from the frontend directory and fix any type errors
- [x] 4.2 Run `npm run lint` and fix any lint issues
