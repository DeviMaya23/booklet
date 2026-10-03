## Why

The Inbox UI has selection and context menus wired up but the "New Artpiece" action is a no-op placeholder. Users can't act on the files they've uploaded — the core workflow of turning inbox files into artpieces is missing.

## What Changes

- New `ArtpieceFormModal` component: a dialog with title (required), notes, artist (dropdown + add), characters (chip multi-select), and a files sub-component.
- New `ArtpieceFilesInput` component: independently-scrolling files region with a sticky dashed-border drop zone at top, file rows (thumbnail, editable name, editable notes, remove button), and upload-within-modal support. Kept as a separate component for future reuse in the "Add to Existing Artpiece" flow.
- Wire up `onNewArtpiece` in `FileContextMenu` (right-click) and in `FilesPage` dropdown, passing selected file IDs into the modal.
- File name/notes edits within the modal are saved on blur via `PUT /files/:id` — no batch on artpiece save.
- Save fires the existing `POST /artpieces` with `file_ids`, which attaches files as part of creation.
- The bruno example for `create_artpiece.bru` is updated to show `file_ids`.

## Capabilities

### New Capabilities

- `web-artpiece-create-modal`: The "Create New Artpiece" modal — form fields, files sub-component behavior (upload, thumbnails, remove), entry points, save flow, and file name/notes save-on-blur semantics.

### Modified Capabilities

- `web-file-inbox`: The "New Artpiece" actions in the selection toolbar dropdown and in the right-click context menu now open the modal instead of being no-ops.

## Impact

- **Frontend**: New components in `frontend/src/features/files/components/` (`ArtpieceFormModal.tsx`, `ArtpieceFilesInput.tsx`). New API hook `useCreateArtpiece`. Wiring changes in `FilesPage.tsx` and `FileInboxGrid.tsx`.
- **Backend**: No changes — `POST /artpieces` already accepts `file_ids` and `PUT /files/:id` already exists.
- **Bruno**: `collection/artpieces/create_artpiece.bru` updated to include `file_ids` in the example body.
- **No breaking changes.**
