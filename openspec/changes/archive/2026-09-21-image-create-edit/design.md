## Context

The Images page already lists and filters images client-side. The backend has all required endpoints: `POST /images` (InitialUpload), `POST /images/:id/complete` (CompleteUpload), `PUT /images/:id` (metadata update), `GET /images/:id` (single image with presigned `image_url`). The character modal (`CharacterFormModal`) establishes the pattern for create/edit modals with two-phase avatar upload, and `TokenInput` + `FolderPicker` establish the pattern for multi-token fields with client-side filtering. The artists feature has `useArtists` and `ArtistFormModal` ready for reuse.

The project uses shadcn on top of `@base-ui/react`. All existing UI components in `src/components/ui/` are thin wrappers over Base-UI primitives. There is no Combobox component yet.

## Goals / Non-Goals

**Goals:**
- Single `ImageFormModal` component that handles both create and edit in the same shape
- Create flow: file picker → preview → two-phase R2 upload on submit
- Edit flow: load current metadata from `GET /images/:id`, submit via `PUT /images/:id`, download full-res image
- Artist field: single-select Combobox with client-side search and inline artist creation via stacked dialog
- Characters field: multi-token input with client-side filtering using existing `TokenInput`
- New `Combobox` UI component in `src/components/ui/` following the existing Base-UI wrapper pattern
- `TokenInput` moved to `src/components/` for shared use

**Non-Goals:**
- Image re-upload / file replacement in edit mode
- Bulk operations
- Server-side search for artists or characters
- Backend changes

## Decisions

### D1: Single modal component for create and edit

`ImageFormModal` accepts an optional `image` prop. When present, it enters edit mode: loads full data via `GET /images/:id` on open, pre-populates fields, shows the download button, hides file upload controls. When absent, it is in create mode: shows the file picker area, requires a file before enabling Save.

**Alternative considered:** Two separate modals (`ImageCreateModal` / `ImageEditModal`). Rejected — the field set is nearly identical and duplication would create maintenance surface. Same pattern as `CharacterFormModal`.

### D2: Edit mode loads from `GET /images/:id`, not from list data

The list endpoint (`GET /images`) omits `image_url` by spec. The full-res presigned URL is required for the download button, so the edit modal must call `GET /images/:id` on open. This also ensures metadata is fresh if the list cache is stale.

The modal stores the fetched `image_url` in local state and shows a loading state while the request is in flight. The rest of the form is disabled until loading completes.

**Alternative considered:** Passing the list item data as the initial values and only fetching `image_url` separately. Rejected — adds complexity for little gain; a single fetch on open is simpler and keeps state in one place.

### D3: Artist field uses a new `Combobox` UI wrapper

`@base-ui/react/combobox` provides `Root`, `Input`, `Trigger`, `Popup`, `Item`, `Empty`, and `useFilter`. A new `src/components/ui/combobox.tsx` wraps these in the same style as the existing `dialog.tsx`/`dropdown-menu.tsx` wrappers — named exports, Tailwind classes, `cn()` for merging.

The artist combobox in `ImageFormModal` is single-select (`multiple` not set). It receives the full artist list from `useArtists()`, filters client-side with Base-UI's `useFilter`, and calls `onValueChange` to update local `artistId` state.

**Alternative considered:** Native `<select>` element. Rejected — no search and inconsistent with design system.

**Alternative considered:** Reusing `TokenInput` for single-select. Rejected — `TokenInput` is built for multi-token entry with free-text commit behaviour; single-select combobox is a different interaction pattern.

### D4: Inline artist creation via stacked dialog

The artist field has a `+` button. Clicking it opens `ArtistFormModal` as a second dialog. Base-UI dialogs render into portals with correct stacking — no z-index conflict with the parent `ImageFormModal`.

After `ArtistFormModal` closes with a successful create, `useCreateArtist`'s `mutateAsync` returns the new `Artist` object. `ImageFormModal` sets that artist as the selected value immediately, without waiting for the artist list to refetch.

`useCreateArtist` already calls `queryClient.invalidateQueries({ queryKey: ARTISTS_QUERY_KEY })` in `onSuccess`, so the artist list stays consistent for future opens.

### D5: Characters field reuses `TokenInput` directly

`TokenInput` already has the shape needed: `suggestions` prop for the dropdown pool, client-side substring filtering, token add/remove, keyboard navigation. `ImageFormModal` passes `useCharacters()` data as suggestions and tracks selected characters in local state as `{ id, name }[]`.

### D6: `TokenInput` moves to `src/components/`

`TokenInput` is now shared between the characters feature (`FolderPicker`) and the images feature (`ImageFormModal`). The correct home is `src/components/TokenInput.tsx`, alongside `ResourceCard.tsx` — these are custom shared components, distinct from the Base-UI wrappers in `src/components/ui/`.

Import paths updated in `FolderPicker.tsx` and `CharacterFormModal.tsx`.

### D7: File upload restricted to jpg/png, required for create

The file input uses `accept="image/jpeg,image/png"`. The Save button in create mode is disabled until a file is picked. On submit, `mime_type` is derived from `file.type`. The character modal's avatar upload logic (`runAvatarUpload`) is the reference implementation: init → PUT to R2 → complete.

If the R2 PUT or CompleteUpload fails during create, the partially created image record is cleaned up with `DELETE /images/:id` (same recovery pattern as character avatar).

### D8: Download image via presigned URL from `GET /images/:id`

In edit mode, the modal stores the `image_url` from the `GET /images/:id` response. The download button creates a temporary `<a>` element with `href=image_url` and `download` attribute and clicks it programmatically. No separate download endpoint is needed.

## Risks / Trade-offs

- **Presigned URL expiry during long sessions** → The URL fetched on modal open is valid for 1 hour. If a user leaves the modal open longer, the download link will expire. Acceptable for now — refreshing the modal re-fetches.
- **Stale artist auto-select after inline create** → After creating an artist via the `+` dialog, we auto-select using the `mutateAsync` return value. If the artist list refetch fails, the displayed name is still correct (we have the object). Low risk.
- **Two-phase upload partial failure** → If CompleteUpload fails, the orphaned image record is deleted. If the DELETE also fails, the record stays orphaned until the stale-upload-purge worker runs. Acceptable — matches existing character avatar behaviour.
