## 1. Shared Infrastructure

- [x] 1.1 Move `TokenInput` from `src/features/characters/components/TokenInput.tsx` to `src/components/TokenInput.tsx`
- [x] 1.2 Update import in `src/features/characters/components/FolderPicker.tsx`
- [x] 1.3 Update import in `src/features/characters/components/CharacterFormModal.tsx`
- [x] 1.4 Create `src/components/ui/combobox.tsx` wrapping `@base-ui/react/combobox` (Root, Input, InputGroup, Trigger, Portal, Popup, Item, Empty, useFilter)

## 2. Image API Hooks

- [x] 2.1 Create `src/features/images/api/useInitImageUpload.ts` — calls `POST /images`, returns `{ id, upload_url, expires_at }`
- [x] 2.2 Create `src/features/images/api/useCompleteImageUpload.ts` — calls `POST /images/:id/complete`
- [x] 2.3 Create `src/features/images/api/useUpdateImage.ts` — calls `PUT /images/:id` with title, notes, artist_id, character_ids

## 3. ImageFormModal — Create Flow

- [x] 3.1 Create `src/features/images/components/ImageFormModal.tsx` with Dialog shell (title "New image" / "Edit image", footer with Cancel + Save)
- [x] 3.2 Implement file picker area: clickable zone, hidden file input (`accept="image/jpeg,image/png"`), local preview on file select, X button to clear
- [x] 3.3 Wire Save disabled state: disabled when no file selected (create mode) or when submission is pending
- [x] 3.4 Implement `handleCreateSubmit`: init upload → PUT to R2 → complete upload; on failure delete the orphaned record and show error toast; on success show success toast, close modal, invalidate `IMAGES_QUERY_KEY`

## 4. ImageFormModal — Edit Flow

- [x] 4.1 Add `image` prop (optional `Image` type) to toggle edit mode
- [x] 4.2 On modal open in edit mode, fetch `GET /images/:id`; disable form while loading; store `image_url` in local state
- [x] 4.3 Pre-populate title, notes, artist, and characters fields from fetched data
- [x] 4.4 Show thumbnail (or placeholder) in the image preview area in edit mode; hide file picker controls
- [x] 4.5 Implement `handleEditSubmit`: call `PUT /images/:id`; on success show toast, close modal, invalidate `IMAGES_QUERY_KEY`
- [x] 4.6 Implement Download image button: programmatic `<a download>` click using stored `image_url`; button only visible in edit mode
- [x] 4.7 Add Delete button (edit mode only) that opens a confirmation AlertDialog; on confirm call `DELETE /images/:id`, close modal, show toast, invalidate `IMAGES_QUERY_KEY`

## 5. Artist Combobox Field

- [x] 5.1 Add artist Combobox to `ImageFormModal` using `useArtists()` as the data source; implement client-side filtering with Base-UI `useFilter`
- [x] 5.2 Wire Combobox `value`/`onValueChange` to local `artistId` state; pass `itemToStringLabel` to display artist name
- [x] 5.3 Add `+` button beside the Combobox; clicking it opens `ArtistFormModal` as a stacked dialog
- [x] 5.4 On `ArtistFormModal` close after a successful create: auto-select the returned artist in the Combobox

## 6. Characters Token Field

- [x] 6.1 Add Characters `TokenInput` to `ImageFormModal` using `useCharacters()` as the suggestions source; exclude already-selected characters from suggestions
- [x] 6.2 Wire `TokenInput` `items`/`onChange` to local `selectedCharacters` state (`{ id, name }[]`)

## 7. Wire Up Pages

- [x] 7.1 Add `onEditClick` prop to `ImagesGrid`; call it when a card is clicked (in addition to existing delete menu)
- [x] 7.2 In `ImagesPage`, add `modalOpen` and `editingImage` state; render `<ImageFormModal>`; wire `+ New` button to open in create mode and card click to open in edit mode with the selected image

## 8. Quality

- [x] 8.1 Run `npm run build` and fix any type errors
- [x] 8.2 Run `npm run lint` and fix any lint errors
