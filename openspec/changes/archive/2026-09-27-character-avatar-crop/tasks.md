## 1. Dependency

- [x] 1.1 Install `react-easy-crop` and its types (`npm install react-easy-crop`)

## 2. Canvas Utility

- [x] 2.1 Create `frontend/src/features/characters/lib/cropCanvas.ts` with `getCroppedFile(sourceUrl, pixelCrop, mimeType): Promise<File>` — load image element, draw cropped region onto offscreen canvas bounded to 1024 × 1024, return File via `toBlob`

## 3. AvatarCropDialog Component

- [x] 3.1 Create `frontend/src/features/characters/components/AvatarCropDialog.tsx` — Dialog wrapping `react-easy-crop` with `aspect={1}` and `cropShape="rect"`, zoom state (0.5–3×), zoom slider, Cancel and Crop buttons
- [x] 3.2 Wire `onCropComplete` callback to capture `PixelCrop`, call `getCroppedFile` on Crop confirm, invoke `onCrop(file)` callback, and revoke source object URL on both confirm and cancel

## 4. CharacterFormModal Updates

- [x] 4.1 Add `dragActive` state and DnD event handlers (`onDragOver`, `onDragLeave`, `onDrop`) to the avatar button — drop calls the same file-handling path as click
- [x] 4.2 Update avatar area empty state to show "Click or drag to upload" helper text alongside the icon; apply ring highlight when `dragActive` is true
- [x] 4.3 Replace direct `setLocalFile` / `setLocalPreviewUrl` in `handleFileChange` (and the new drop handler) with `stageFileForCrop` — stages the source file and opens `AvatarCropDialog` (no size gate; the crop editor handles any input)
- [x] 4.4 Add `cropSourceFile` and `cropDialogOpen` state; on crop confirm set `localFile` + `localPreviewUrl` from the cropped File and close the crop dialog; on crop cancel close the dialog with no state change
- [x] 4.5 Render `AvatarCropDialog` inside the modal, driven by `cropDialogOpen` and `cropSourceFile`

## 5. Verification

- [x] 5.1 Run `npm run build` and fix any TypeScript or build errors
- [x] 5.2 Run `npm run lint` and fix any lint issues
