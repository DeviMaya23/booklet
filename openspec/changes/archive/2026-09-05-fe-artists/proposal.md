## Why

The artists page (`/app/artists`) exists as a stub returning null. Users have no UI to manage artists — creating, editing, linking, or deleting them — despite the backend being fully implemented.

## What Changes

- Implement the `/app/artists` page with a searchable list of artists
- Each artist row displays: a bullet, the artist name, a copy-link button, and an edit button
- A shared create/edit modal handles both create and edit modes, with delete available in edit mode only
- Client-side search filters artists by name (matching the images/characters pattern)

## Capabilities

### New Capabilities

- `web-artists`: UI for listing, searching, creating, editing, and deleting artists in the web app

### Modified Capabilities

## Impact

- `frontend/src/pages/ArtistsPage.tsx` — implement (currently returns null)
- New feature module `frontend/src/features/artists/` — api hooks and list/modal components
- No backend changes required — all endpoints are already implemented
