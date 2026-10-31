# Proposal

## Why

In the artpiece detail view's edit mode, the Artist dropdown has no way to create a new artist. The artpiece create modal and the commission form modal both have a "+" button next to the dropdown, so a user who needs a new artist while editing has to leave for the Artists page, create it, and come back. Edit mode should offer the same pathway.

## What Changes

- Add a "+" button next to the Artist combobox in `ArtpieceDetailView` edit mode, matching the existing pattern in `ArtpieceFormModal` and `CommissionFormModal`.
- Clicking it opens the existing `ArtistFormModal`; on successful creation the new artist is selected in the edit form (`onCreated` sets the edit artist).
- Creating an artist does not save the artpiece; the selection counts as a pending change and appears in the unsaved-changes caption like any other artist change.
- No new component or abstraction: the wiring is copied locally, as the two existing call sites do.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `web-artpiece-detail`: edit mode's Artist field gains an add-artist action that opens the artist-creation modal and selects the created artist.

## Impact

- Frontend only: `frontend/src/features/artpieces/components/ArtpieceDetailView.tsx` and its test file `ArtpieceDetailView.test.tsx`.
- No backend, API, or dependency changes; no new endpoints, so no Bruno file.
- Reuses `ArtistFormModal` and `useCreateArtist`, which already invalidates `ARTISTS_QUERY_KEY` so the new artist appears in the combobox list.
