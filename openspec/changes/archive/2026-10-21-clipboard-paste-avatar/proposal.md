# Proposal

## Why

The character avatar upload currently requires either clicking the avatar button (file picker) or dragging a file into the drop zone. Users who have screenshot-copied an image have no way to paste it directly — they must first save it as a file, then upload it. Clipboard paste is a natural, low-friction path that costs nothing downstream since the crop flow already handles any `File` object.

## What Changes

- `CharacterFormModal` gains a `paste` event listener on the dialog container that intercepts `image/*` clipboard items and routes them into the existing `stageFileForCrop` path.
- Paste is suppressed when the event originates from a focused text input or textarea so that normal text paste is unaffected.
- No changes to `AvatarCropDialog`, `getCroppedFile`, or any backend endpoint.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `web-character-avatar-crop`: The crop flow currently accepts images from a file picker or drag-and-drop. A new input path — clipboard paste — routes images into the same `stageFileForCrop` function. The requirement for what the crop dialog does is unchanged; only the trigger surface grows.

## Impact

- **Frontend only** — `CharacterFormModal.tsx` is the sole file changed.
- No new dependencies.
- No API or backend changes.
- No changes to `AvatarCropDialog` or `cropCanvas`.
