# Proposal

## Why

Saved filters exist in the backend but are completely unwired. Users have no way to save the current artpiece gallery filter state to the dashboard. This is Part 1 of wiring: add the ability to create a saved filter directly from the artpiece list's filter popover.

## What Changes

- **ArtistCharacterFilter component**: gains an optional `footerAction?: ReactNode` slot rendered in the popover footer alongside "Clear all", so callers can inject context-specific actions without modifying the shared component.
- **Artpiece gallery filter persistence**: filter state (artist, characters, match mode) is synced to URL search params so navigating away and back restores the previous filter.
- **"Save to dashboard" button**: appears in the filter popover footer when at least one filter is active; opens a save modal.
- **Save filter modal** (`SaveFilterModal`): shows a read-only filter summary chip, a name text input, a thumbnail grid (from the filtered artpieces), a dashboard preview, and Save/Cancel. Save is disabled when the filtered result is empty. After a successful save, the modal closes and a success toast fires; the popover footer shows "Saved as `<name>`" in place of the button as long as the filter hasn't changed.
- **Backend create request change**: `POST /saved_filters` accepts `thumbnail_artpiece_id` (UUID) instead of `thumbnail_r2_path`. The handler resolves the artpiece's cover thumbnail R2 path and copies it into the saved filter record. The handler gains a presigner and artpiece repository dependency.
- **Backend duplicate name**: `POST /saved_filters` returns `422` when the user already has a saved filter with the same name.
- **Backend response change**: `GET /saved_filters`, `GET /saved_filters/:id`, `PATCH /saved_filters/:id`, and `POST /saved_filters` responses replace `thumbnail_r2_path` with `thumbnail_url` (a pre-signed URL generated at response time, null if no thumbnail).
- **DB migration**: unique constraint on `(user_id, name)` in `saved_filters`.
- **New FE API hooks**: `useSavedFilters` (query) and `useCreateSavedFilter` (mutation).
- **Bruno files**: update create file to reflect new request shape; add a list file that now shows `thumbnail_url`.

## Capabilities

### New Capabilities

- `web-save-filter-modal`: the "Save to dashboard" modal component — filter summary display, name input, thumbnail picker from filtered artpieces, dashboard preview tile, and save action.

### Modified Capabilities

- `saved-filters`: create request replaces `thumbnail_artpiece_id` for `thumbnail_r2_path`; duplicate name returns 422; all responses replace raw `thumbnail_r2_path` with signed `thumbnail_url`.
- `web-artpieces-gallery`: filter state is persisted to URL params; filter popover gains a "Save to dashboard" footer action wired through the `footerAction` slot.
- `web-artist-character-filter`: component gains an optional `footerAction?: ReactNode` slot in the popover footer.

## Impact

- **Backend**: `saved_filter_handler.go` — constructor gains `ArtpieceRepository` and `Presigner`; create handler resolves artpiece thumbnail; list/get/patch handlers sign the thumbnail path; all handler tests updated.
- **Backend**: new DB migration for the unique constraint.
- **Frontend**: `ArtistCharacterFilter.tsx` and its test — new optional prop.
- **Frontend**: `useArtpiecesFilter.ts` — filter state moves to URL params via `useSearchParams`.
- **Frontend**: `ArtpiecesPage.tsx` — wires `footerAction` slot; manages "saved as" ephemeral state.
- **Frontend**: new `SaveFilterModal.tsx` and `useSavedFilters.ts` / `useCreateSavedFilter.ts`.
- **Bruno**: `collection/saved_filters/create_saved_filter.bru` updated; responses now show `thumbnail_url`.
- **No dashboard changes** in this proposal (Part 2).
- `CommissionFormModal` and `AddToExistingArtpieceModal` are unaffected — they pass no `footerAction`.
