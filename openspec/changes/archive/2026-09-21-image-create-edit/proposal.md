## Why

The Images page has a working list and client-side search, but the `+ New` button is a no-op and cards have no edit entry point. Users cannot create images or update image metadata from the UI.

## What Changes

- Add `ImageFormModal` covering create (two-phase upload) and edit (metadata update) flows
- Wire the `+ New` button and card clicks in `ImagesPage` and `ImagesGrid` to the modal
- Add a `Combobox` UI component wrapping `@base-ui/react/combobox` for single-select with search
- Move `TokenInput` from `features/characters/components/` to `src/components/` for reuse across features
- Add FE API hooks for image upload init, upload complete, and image update

## Capabilities

### New Capabilities

- `web-images-create`: User can create a new image record by picking a jpg/png file, filling in title/notes/artist/characters, and submitting — triggers the two-phase R2 upload flow
- `web-images-edit`: User can open an existing image card to edit its title, notes, artist, and characters, and download the full-res image

### Modified Capabilities

- `web-images`: Add requirement that the `+ New` button opens the create modal and each image card opens the edit modal

## Impact

**Frontend only — no backend changes required.** All endpoints (`POST /images`, `POST /images/:id/complete`, `PUT /images/:id`, `GET /images/:id`) are already implemented.

New files:
- `src/components/TokenInput.tsx` (moved from `features/characters/components/`)
- `src/components/ui/combobox.tsx` (new Base-UI Combobox wrapper)
- `src/features/images/api/useInitImageUpload.ts`
- `src/features/images/api/useCompleteImageUpload.ts`
- `src/features/images/api/useUpdateImage.ts`
- `src/features/images/components/ImageFormModal.tsx`

Modified files:
- `src/features/images/components/ImagesGrid.tsx`
- `src/pages/ImagesPage.tsx`
- `src/features/characters/components/FolderPicker.tsx` (update import)
- `src/features/characters/components/CharacterFormModal.tsx` (update import)
