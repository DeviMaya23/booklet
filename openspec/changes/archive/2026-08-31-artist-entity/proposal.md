## Why

Images currently store artist attribution as free-text fields (`artist_name`, `artist_link`) that must be re-typed for every upload. Extracting artist into a first-class entity lets users register an artist once and reference it across images.

## What Changes

- New `Artist` entity with fields: `id`, `user_id`, `name`, `notes`, `artist_link`; name is unique per user at the DB level
- New CRUD endpoints for artists (`POST /artists`, `GET /artists`, `GET /artists/:id`, `PATCH /artists/:id`, `DELETE /artists/:id`)
- `Image` replaces `artist_name` + `artist_link` columns with a nullable `artist_id` FK to `artists` — **BREAKING** for any client reading those fields
- `PendingUpload` similarly replaces `artist_name` + `artist_link` with nullable `artist_id` FK
- `POST /images` (initial upload) now accepts `artist_id` instead of inline artist fields
- `PATCH /images/:id` now accepts `artist_id` instead of inline artist fields
- `CompleteUpload` validates `artist_id` still exists and belongs to the user; silently sets `artist_id = null` on the image if not
- Deleting an artist sets `artist_id = null` on all referencing images and pending uploads (`ON DELETE SET NULL`)
- Existing `artist_name`/`artist_link` data in the DB is dropped by the migration (no backfill)

## Capabilities

### New Capabilities

- `artist-management`: CRUD for the artist entity scoped per user — create, list, get by ID, update, delete

### Modified Capabilities

- `image-management`: image response shape changes (drops `artist_name`/`artist_link`, adds `artist_id` + embedded artist name); update request changes accordingly
- `image-upload`: initial upload request drops inline artist fields, accepts `artist_id`; complete upload gains artist existence check

## Impact

- **DB migration**: new `artists` table; `images` and `pending_uploads` columns altered
- **Backend**: new domain entity, repository, usecase, handler; modifications to image and upload layers across all files
- **API**: breaking change to `GET /images`, `GET /images/:id`, `PATCH /images/:id`, `POST /images` response/request shapes
- **Clients**: any consumer of artist fields on image endpoints must update to `artist_id`-based model
