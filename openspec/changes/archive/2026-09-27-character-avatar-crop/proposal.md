## Why

Character avatars are always displayed as squares, but the existing avatar picker accepts any image as-is — resulting in poorly framed previews when the subject isn't centred or the source is portrait/landscape. Adding a crop step ensures every avatar looks intentional before it's uploaded.

## What Changes

- The avatar area in `CharacterFormModal` gains helper text ("Click or drag to upload") and drag-and-drop support, consistent with the image modal
- Any file pick (click or drop, in both create and edit mode) is intercepted before `localFile` is set — a crop dialog opens instead
- The new `AvatarCropDialog` lets the user pan, zoom, and confirm a square crop; only after confirmation does the cropped `File` flow into the existing upload logic
- A raw file size gate (>10 MB) rejects obviously large inputs before the crop dialog opens
- The crop output is always square JPEG, bounded to 1024 × 1024 px, which naturally keeps upload size small
- The avatar delete flow is unchanged

## Capabilities

### New Capabilities

- `web-character-avatar-crop`: Square crop dialog — zoom slider, pan, confirm/cancel — that intercepts the avatar file pick and produces a cropped File for the existing upload flow

### Modified Capabilities

- `web-characters`: Avatar picker requirement changes — adds helper text, drag-and-drop, and mandates all file picks go through the crop step before preview/upload

## Impact

- **Frontend only** — no backend or API changes
- New npm dependency: `react-easy-crop`
- New component: `frontend/src/features/characters/components/AvatarCropDialog.tsx`
- New utility: `frontend/src/features/characters/lib/cropCanvas.ts`
- Modified: `frontend/src/features/characters/components/CharacterFormModal.tsx`
