## Why

The app currently stores artwork as flat images with metadata attached directly to each image. This doesn't support artworks with multiple file versions (e.g. a PSD source alongside a rendered PNG, or an MP4 alongside a thumbnail image), and conflates file concerns with artwork metadata. Introducing artpieces as the canonical browsable unit — grouping one or more files and carrying the artwork's metadata — unblocks non-image file types and gives the gallery a stable identity to build on.

## What Changes

- New `artpieces` table: the browsable unit carrying title, artist, notes, characters, and a cover file
- New `files` table: generic file storage derived from `images`, with `artpiece_id` linking to its owner
- New `image_metadata` side table: width/height for image-type files only
- New `artpiece_characters` join table: many-to-many between artpieces and characters
- New `pending_file_uploads` table: tracks in-flight file uploads (separate from existing `pending_uploads`)
- New CRUD endpoints for artpieces: create, get, list (with filters), update, delete
- New file attach/detach endpoints: attach a file to an artpiece, detach it, set cover
- New file upload endpoints: init upload (presigned PUT to R2), complete upload (create file row, enqueue thumbnail for images)
- New `generate_file_thumbnail` background worker (parallel to existing `generate_thumbnail` for images)
- Move `mimeTypeToExt` helper from `upload_usecase.go` to `pkg/mime` for shared use
- Existing `images` tables, endpoints, and frontend remain in place and unchanged

## Capabilities

### New Capabilities

- `artpiece-management`: CRUD for artpieces — create, get, list (with artist/character filters), update, delete
- `artpiece-files`: Attach/detach files to artpieces; set cover; cover auto-selection and cleanup rules
- `file-upload`: Upload flow for files — presigned PUT init, complete with R2 storage, thumbnail generation for image mime types, `image_metadata` row for images

### Modified Capabilities

## Impact

- **New migrations**: `artpieces`, `files`, `artpiece_characters`, `image_metadata`, `pending_file_uploads` tables; circular FK between `artpieces.cover_file_id` and `files.artpiece_id` resolved via deferred ALTER
- **New backend packages/files**: `domain.Artpiece`, `domain.File`, `domain.ImageMetadata`, `domain.PendingFileUpload`; artpiece repository, usecase, handler; file upload usecase, handler; `worker.GenerateFileThumbnailWorker`; `pkg/mime` package
- **New Bruno collection files**: one per endpoint under `collection/artpieces/` and `collection/files/`
- **No changes** to `images`, `image_characters`, `pending_uploads`, or any existing endpoints/frontend
- **Existing data**: disposable; no migration of existing image rows into the new schema
