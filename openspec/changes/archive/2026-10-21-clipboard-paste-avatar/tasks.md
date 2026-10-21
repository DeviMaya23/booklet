# Tasks

## 1. Clipboard paste handler

- [x] 1.1 Add an `onPaste` handler to `<DialogContent>` in `CharacterFormModal.tsx` that: checks if `event.target` is `HTMLInputElement` or `HTMLTextAreaElement` (return early if so); finds the first `image/*` item in `event.clipboardData.items`; calls `stageFileForCrop(file)` with the result of `getAsFile()`. Verify by pasting a screenshot into the open modal (not in a text field) and confirming the crop dialog opens.

- [x] 1.2 Verify that pressing Ctrl+V while the name input is focused pastes text normally and does not open the crop dialog. Verify that pressing Ctrl+V while the notes textarea is focused behaves the same way.

- [x] 1.3 Verify that pasting an image when an avatar already exists opens the crop dialog and, after cropping, replaces the avatar preview with the new image.

- [x] 1.4 Run `npm run build` and `npm run lint` in `frontend/` and fix any errors or warnings that arise.

## Workflow follow-up

- Archive the change once review is complete.
