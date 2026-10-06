## Why

When users upload files, the filename (including its extension, e.g. `pic.jpg`) is stored as the file's `name`. This creates a double-extension problem when downloads reconstruct the filename by appending an extension derived from MIME type — resulting in `pic.jpg.jpg`. Additionally, no validation prevents users from entering names with characters that are illegal in filenames on Windows and macOS.

## What Changes

- When a file is added via any upload entry point (DnD or file picker), the extension is automatically stripped from `file.name` before it is sent to the backend as `name` — so `pic.jpg` becomes `pic`, and the user's name input is pre-filled with the clean name.
- File name inputs across all upload surfaces validate on blur against a forbidden-character set (`/ \ : * ? " < > |`), show an inline error, and block the save call until the name is valid.

## Capabilities

### New Capabilities

- `web-file-name-validation`: Inline validation on file name inputs — forbidden characters, blank name, and max-length (255) rules, applied consistently across all file upload and edit surfaces.

### Modified Capabilities

- `web-file-inbox`: Inbox DnD upload strips extension from filename before sending `name` to initiate-upload.
- `web-file-inbox-upload-modal`: Upload modal strips extension from filename before pre-filling the name input and before sending `name` to initiate-upload. Name input validates on blur with the shared validation rules.
- `web-artpiece-create-modal`: `ArtpieceFilesInput` component strips extension before pre-filling the name input and before sending `name` to initiate-upload. Name input validates on blur.
- `web-artpiece-detail`: Edit-mode DnD/file-picker upload in the detail view strips extension before sending `name` to initiate-upload.

## Impact

- Frontend only — no backend changes, no migrations.
- Four upload call sites updated: `AddFilesModal.tsx`, `ArtpieceFilesInput.tsx`, `FileInboxGrid.tsx`, `ArtpieceDetailView.tsx`.
- A shared `stripFileExtension` utility function introduced in `frontend/src/lib/` (or equivalent shared utils location) to avoid duplication.
- Name validation logic can be a shared `validateFileName` utility used by all name inputs.
