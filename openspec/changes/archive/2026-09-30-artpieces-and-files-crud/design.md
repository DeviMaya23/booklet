## Context

The app's current image flow is a single-concern pipeline: each image row carries both the file reference (R2 path, mime type) and the artwork metadata (title, artist, characters). This works for single-file artworks but collapses when an artwork has multiple file representations (e.g. a PSD source + rendered PNG, or an MP4 + cover image).

The new model separates these concerns: `artpieces` own metadata, `files` own storage. The existing `images` table and all endpoints remain untouched; the new tables coexist with them until the images flow is retired in a separate proposal.

## Goals / Non-Goals

**Goals:**
- Introduce `artpieces` as the browsable unit with metadata (title, artist, characters, notes, cover)
- Introduce `files` as the generic storage unit, derived from `images` but without artwork metadata
- Ship CRUD for artpieces and an upload flow for files that mirrors the existing image upload flow
- Generate thumbnails and `image_metadata` rows for image-type files at upload completion
- Enforce ownership scoping: user can only reference their own artists, characters, and files

**Non-Goals:**
- Deleting, migrating, or modifying the `images` table, existing endpoints, or frontend
- Migrating existing image rows into the new schema
- File type validation (no mime type whitelist)
- Frontend implementation for artpieces

## Decisions

### D1: Separate `pending_file_uploads` table

The existing `pending_uploads` table carries `title`, `artist_id`, and `character_ids` — fields that belong to artpieces, not files. Reusing it for file uploads would leave those columns permanently null for files. A separate `pending_file_uploads` table is introduced with only the fields a file upload needs: `mime_type`, `artpiece_id` (optional), and `notes`.

_Alternative_: Reuse `pending_uploads` with null artwork fields. Rejected because it blurs the separation of concerns and leaves the schema misleading long-term.

### D2: Circular FK resolved with deferred ALTER TABLE

`artpieces.cover_file_id` → `files` and `files.artpiece_id` → `artpieces` reference each other. Creating both tables in the same migration with both FKs is not possible. Resolution:

1. Migration N: create `artpieces` **without** `cover_file_id`
2. Migration N+1: create `files` with `artpiece_id FK → artpieces`
3. Migration N+2: `ALTER TABLE artpieces ADD COLUMN cover_file_id uuid REFERENCES files ON DELETE SET NULL`

`cover_file_id` is nullable, so no data backfill is needed.

### D3: Cover integrity enforced in application code

A composite FK (`cover_file_id` must reference a file whose `artpiece_id` is the same artpiece) cannot be expressed as a standard Postgres FK. Enforcement lives in the usecase layer:

- **Set cover**: verify `file.artpiece_id == artpiece.id` and `file.user_id == user.id` before writing
- **Detach file**: if the detached file was the cover, run cover reassignment (image-first, then any file, then null)
- **Upload complete**: if file is attached to an artpiece with no cover, auto-set cover (image-first)

_Alternative_: Postgres composite FK via a unique constraint on `(id, artpiece_id)` on `files` plus a FK on `(cover_file_id, id)` on `artpieces`. Possible, but requires non-standard GORM workarounds and adds migration complexity. Application enforcement is simpler and consistent with existing patterns in the codebase.

### D4: Cover auto-selection logic

When a file is attached to an artpiece that has no cover:
1. If the new file is an image (mime type starts with `image/`) → set as cover
2. If no image files exist in the artpiece → set this file as cover regardless of type

When a file leaves an artpiece (detach) and it was the cover:
1. Find another file in the artpiece where mime type starts with `image/` → set as cover
2. No image files remain → set first remaining file as cover
3. No files remain → clear cover (NULL)

### D5: Thumbnail generation and `image_metadata` gated on mime type

At file upload completion, two operations are conditional on `strings.HasPrefix(mimeType, "image/")`:
1. Enqueue `generate_file_thumbnail` River job
2. Insert `image_metadata` row (width/height)

No mime type whitelist is enforced. If a client claims an image mime type but uploads a non-image file, the thumbnail job will fail in River (retried per River's default policy) and log an error. This is acceptable; there is no frontend yet and data is disposable.

### D6: `mimeTypeToExt` moved to `pkg/mime`

The helper is currently private in `upload_usecase.go`. Both the existing upload usecase and the new file upload usecase need it. It has no app-specific knowledge, so it belongs in `pkg/mime`.

### D7: Separate `GenerateFileThumbnailWorker`

The existing `GenerateThumbnailWorker` reads from `imageRepo` and writes back to `imageRepo.UpdateThumbnailPath`. Rather than coupling the new file thumbnail flow to the images repo, a separate `GenerateFileThumbnailWorker` reads from `fileRepo` and writes to `fileRepo.UpdateThumbnailPath`. This keeps the workers independent.

### D8: List query uses GORM preloads

The artpiece list endpoint preloads `CoverFile`, `Artist`, and `Characters`. This mirrors how the image list preloads `Artist` and `Characters`. No raw SQL joins.

## Risks / Trade-offs

- **Mime type spoofing on thumbnail jobs** → River will retry and fail for non-decodable files; fails silently in the background. Acceptable given no FE and disposable data.
- **Cover integrity in concurrent writes** → two simultaneous detach/attach operations on the same artpiece could race on cover assignment. No transaction wrapping of the cover-check + cover-write today. Low risk at current scale; worth noting for future.
- **Circular FK migration ordering** → a failed migration at step N+2 leaves artpieces without the `cover_file_id` column. The app would boot but cover writes would fail. Rollback: `ALTER TABLE artpieces DROP COLUMN cover_file_id`.

## Migration Plan

Migrations run automatically at server boot via golang-migrate. Order:

1. `000015_create_artpieces_no_cover.up.sql` — artpieces table without cover_file_id
2. `000016_create_files.up.sql` — files table with artpiece_id FK
3. `000017_create_artpiece_characters.up.sql` — join table
4. `000018_create_image_metadata.up.sql` — side table for image dimensions
5. `000019_add_cover_to_artpieces.up.sql` — ALTER TABLE adds cover_file_id FK
6. `000020_create_pending_file_uploads.up.sql` — pending upload tracking for files

Rollback: each migration has a corresponding `.down.sql`. The server does not run down migrations automatically; they are applied manually if needed.
