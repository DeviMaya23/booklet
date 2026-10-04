## Why

The "+" floating button on the Inbox page is currently a no-op. Users need a first-party upload entry point for when they want to name or annotate files as they upload them — distinct from drag-and-drop onto the grid, which auto-uploads with no annotation opportunity.

## What Changes

- The floating "+" button on `/app/files` opens a new "Add files to dump" modal
- The modal provides a drag-and-drop drop zone and file picker to upload multiple files in parallel
- Each uploaded file appears as a row with a thumbnail, editable name field, and editable notes field
- Name and notes save automatically on blur (no global Save/Cancel — uploads commit immediately, same as everywhere else)
- Upload failure: row is removed from the list and an error toast is shown naming the file
- Field-save failure: the field shows a red border and an inline error message; clears on successful retry
- No X button while a file is uploading; X appears only after upload completes
- Footer contains only a Close button (dismissive only — nothing to commit or discard)

## Capabilities

### New Capabilities

- `web-file-inbox-upload-modal`: The "Add files to dump" modal — upload entry point behind the "+" button, with per-row auto-save and inline field-error feedback

### Modified Capabilities

- `web-file-inbox`: The "+" button now opens the upload modal instead of being a no-op

## Impact

- `frontend/src/pages/FilesPage.tsx`: wire the "+" button to open the new modal
- New component: `frontend/src/features/files/components/AddFilesModal.tsx`
- Reuses existing hooks: `useInitFileUpload`, `useCompleteFileUpload`, `useUpdateFile`, `useFiles` (for thumbnail polling)
- No backend changes required — all existing endpoints are sufficient
