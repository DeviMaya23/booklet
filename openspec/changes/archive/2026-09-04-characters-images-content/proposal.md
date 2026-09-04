## Why

The app shell routes for `/app/characters` and `/app/images` exist but render nothing. Users need browsable, searchable content areas for their two primary resource types.

## What Changes

- Add a resource card component (`ResourceCard`) shared by both pages: square aspect-ratio card showing a thumbnail/avatar image (or a "no image" placeholder icon), a resource name/title in the bottom-left overlay, and a `(...)` menu in the top-right with a Delete action.
- Add the Characters content page: search bar (client-side filter), `+ New` button (no-op), and a responsive grid of character cards fetched from `GET /characters`.
- Add the Images content page: same layout, grid of image cards fetched from `GET /images`.
- Add a shared Delete confirmation dialog used by both pages.

## Capabilities

### New Capabilities

- `web-characters`: Characters content area at `/app/characters` — searchable grid of character cards with per-card delete.
- `web-images`: Images content area at `/app/images` — searchable grid of image cards with per-card delete.

### Modified Capabilities

## Impact

- Frontend only — no backend changes.
- New files: `ResourceCard`, feature-level API hooks and grid components for characters and images.
- Installs no new dependencies — uses `@tanstack/react-query`, `lucide-react`, `sonner`, and shadcn/ui primitives already present.
- Delete endpoints (`DELETE /characters/:id`, `DELETE /images/:id`) are existing; this is their first UI surface.
