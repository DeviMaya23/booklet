## Why

There is no dedicated place to browse artpieces — they can only be created from the file inbox, with no way to view, search, or filter the full collection. Adding an Artpieces gallery screen gives the artpiece the same first-class navigation presence as Characters, Images, and Artists.

## What Changes

- Add an **Artpieces** sidebar nav item between Inbox and Characters at `/app/artpieces`
- New gallery screen with:
  - Search by title (client-side)
  - Sort dropdown: alphabetical, newest first (client-side)
  - Filter popover: single-artist select (Combobox) + character multi-select with All/Any toggle (client-side)
  - Card grid reusing the existing `ResourceCard` component, showing thumbnail, title, and artist sublabel
  - **+ New Artpiece** button opening the existing `ArtpieceFormModal` with an empty file list
- Add optional `sublabel` prop to `ResourceCard` to support the artist line beneath the title

## Capabilities

### New Capabilities
- `web-artpieces-gallery`: Artpieces gallery page — search, sort, filter popover, card grid, and new artpiece entry point

### Modified Capabilities
- `web-app-shell`: Sidebar nav requirement changes to include Artpieces between Inbox and Characters

## Impact

- **Frontend**: New `ArtpiecesPage.tsx`; new `ArtpiecesGrid` component; `ResourceCard` gains optional `sublabel`; `AppSidebar` and `App.tsx` updated; `useArtpieces` hook already exists
- **Backend**: No changes — `GET /artpieces` already returns `thumbnail_url`, `artist_name`, and `characters[]`; all filtering is client-side
- **Specs**: No changes to artpiece-management, web-artpiece-create-modal, or file-related specs
