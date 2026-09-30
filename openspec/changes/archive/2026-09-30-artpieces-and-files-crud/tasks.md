## 1. Infrastructure: pkg/mime + migrations

- [x] 1.1 Create `pkg/mime` package: move `mimeTypeToExt` from `internal/usecase/upload_usecase.go`, add a `IsImage(mimeType string) bool` helper (`strings.HasPrefix(mimeType, "image/")`)
- [x] 1.2 Update `upload_usecase.go` to call `mime.MimeTypeToExt` instead of the local function
- [x] 1.3 Write migration `000015_create_artpieces_no_cover.up/down.sql` — `artpieces` table without `cover_file_id`
- [x] 1.4 Write migration `000016_create_files.up/down.sql` — `files` table with `artpiece_id FK → artpieces ON DELETE SET NULL`, index on `user_id`
- [x] 1.5 Write migration `000017_create_artpiece_characters.up/down.sql` — join table with cascade
- [x] 1.6 Write migration `000018_create_image_metadata.up/down.sql` — `file_id PK FK → files ON DELETE CASCADE`, `width`, `height`
- [x] 1.7 Write migration `000019_add_cover_to_artpieces.up/down.sql` — `ALTER TABLE artpieces ADD COLUMN cover_file_id uuid REFERENCES files ON DELETE SET NULL`
- [x] 1.8 Write migration `000020_create_pending_file_uploads.up/down.sql` — `id`, `user_id`, `r2_key`, `mime_type`, `artpiece_id` (nullable), `notes` (nullable), `created_at`

## 2. Domain types

- [x] 2.1 Add `domain.File` struct (GORM model): `ID`, `UserID`, `ArtpieceID` (*uuid.UUID), `FileR2Path`, `MimeType`, `ThumbnailR2Path` (*string), `Notes` (*string), `CreatedAt`, `UpdatedAt`
- [x] 2.2 Add `domain.ImageMetadata` struct: `FileID` (PK), `Width`, `Height`; add `ImageMetadata *ImageMetadata` association on `domain.File`
- [x] 2.3 Add `domain.Artpiece` struct: `ID`, `UserID`, `Title` (*string), `ArtistID` (*uuid.UUID), `CoverFileID` (*uuid.UUID), `Notes` (*string), `CreatedAt`, `UpdatedAt`; associations: `Artist *Artist`, `CoverFile *File`, `Characters []Character` (many2many:`artpiece_characters`)
- [x] 2.4 Add `domain.PendingFileUpload` struct: `ID`, `UserID`, `R2Key`, `MimeType`, `ArtpieceID` (*uuid.UUID), `Notes` (*string), `CreatedAt`

## 3. File upload flow

- [x] 3.1 Create `internal/usecase/file_upload_repository.go` — define `FileUploadRepository`, `FileUploadArtpieceRepository`, `FileUploadFileRepository` interfaces
- [x] 3.2 Create `internal/usecase/file_upload_usecase.go` — `InitiateUpload` (validate artpiece ownership, presign PUT, create pending record) and `CompleteUpload` (delete pending, create file row, insert image_metadata if image, enqueue thumbnail job if image, trigger cover auto-select)
- [x] 3.3 Write unit tests for `file_upload_usecase.go`: cover auto-select on first image attach; no thumbnail enqueued for non-image type
- [x] 3.4 Create `internal/repository/file_upload_repository.go` — implement `PendingFileUploadRepository` (Create, GetByID, Delete, ListStale)
- [x] 3.5 Create `internal/repository/file_repository.go` — implement file repo methods needed by upload: `Create`, `GetByIDAndUserID`, `UpdateThumbnailPath`
- [x] 3.6 Write integration tests for file upload repository: create file, get with wrong user returns not found, update thumbnail path
- [x] 3.7 Create `internal/handler/file_upload_handler.go` — `InitiateUpload` (POST /files) and `CompleteUpload` (POST /files/:id/complete)
- [x] 3.8 Write handler unit tests for `file_upload_handler.go`: happy path (assert response body), artpiece not owned → 422, pending not found → 404, invalid UUID path param → 400

## 4. Thumbnail worker for files

- [x] 4.1 Create `internal/worker/generate_file_thumbnail.go` — `GenerateFileThumbnailArgs` and `GenerateFileThumbnailWorker`; reads file from `fileRepo`, downloads from R2, resizes to 600×600, encodes JPEG 85%, uploads to `users/{user_id}/thumbnails/{file_id}.jpg`, calls `fileRepo.UpdateThumbnailPath`
- [x] 4.2 Add `fileRepo.GetByIDForWorker` to file repository and its worker-facing interface
- [x] 4.3 Add stale pending file upload purge worker (reuse pattern from `purge_upload.go`): periodic River job that calls `fileUploadUsecase.CleanupStaleUploads`

## 5. Artpiece CRUD

- [x] 5.1 Create `internal/usecase/artpiece_repository.go` — define `ArtpieceRepository` interface (Create, GetByID, List, Update, Delete)
- [x] 5.2 Create `internal/usecase/artpiece_usecase.go` — implement Create, GetByID, List, Update, Delete; validate artist and character ownership on create/update
- [x] 5.3 Write unit tests for `artpiece_usecase.go`: artist not owned → ErrArtistNotOwned; character not owned → ErrCharacterNotOwned; clear artist (null artist_id on update)
- [x] 5.4 Create `internal/repository/artpiece_repository.go` — GORM implementation; List preloads `CoverFile`, `Artist`, `Characters`; filter by artist_ids and character_ids
- [x] 5.5 Write integration tests for artpiece repository: query correctness (preloads), ownership isolation (wrong user → not found), filter by artist, filter by character, delete does not cascade to files
- [x] 5.6 Create `internal/handler/artpiece_handler.go` — Create (POST), GetByID (GET /:id), List (GET), Update (PUT /:id), Delete (DELETE /:id)
- [x] 5.7 Write handler unit tests for `artpiece_handler.go`: happy path per endpoint (status + body), not found → 404, artist not owned → 422, character not owned → 422, invalid UUID → 400, generic error → 500

## 6. Artpiece file attach/detach/cover

- [x] 6.1 Extend `artpiece_usecase.go` with `AttachFile`, `DetachFile`, `SetCover` methods; enforce ownership of file and artpiece; run cover auto-select/reassignment logic
- [x] 6.2 Write unit tests for attach/detach/cover usecase logic: auto-set cover on first image; auto-set cover on first non-image when no images; cover unchanged when artpiece already has cover; cover reassigned to image on detach; cover cleared when no files remain; set cover rejects file from wrong artpiece
- [x] 6.3 Extend `artpiece_repository.go` with methods needed by attach/detach: `UpdateCover`, `GetFilesForArtpiece`
- [x] 6.4 Extend `file_repository.go` with `UpdateArtpieceID` (for attach/detach)
- [x] 6.5 Add attach/detach/cover endpoints to `artpiece_handler.go`: POST /artpieces/:id/files/:file_id, DELETE /artpieces/:id/files/:file_id, PUT /artpieces/:id/cover
- [x] 6.6 Write handler unit tests for attach/detach/cover: happy path, file not owned → 422, artpiece not found → 404, file not attached to artpiece → 422

## 7. Wire into server

- [x] 7.1 Register `domain.File`, `domain.Artpiece`, `domain.PendingFileUpload` — confirm GORM automigrate is not in use (it isn't; migrations handle schema)
- [x] 7.2 Instantiate all new repositories, usecases, and handlers in `cmd/server/main.go` (`initApp`)
- [x] 7.3 Register new routes in `initApp`: `/files`, `/files/:id/complete`, `/artpieces` and sub-routes
- [x] 7.4 Add `GenerateFileThumbnailWorker` and stale pending file purge periodic job to River workers

## 8. Bruno collection

- [x] 8.1 Create `collection/artpieces/create_artpiece.bru`
- [x] 8.2 Create `collection/artpieces/get_artpiece.bru`
- [x] 8.3 Create `collection/artpieces/list_artpieces.bru`
- [x] 8.4 Create `collection/artpieces/update_artpiece.bru`
- [x] 8.5 Create `collection/artpieces/delete_artpiece.bru`
- [x] 8.6 Create `collection/artpieces/attach_file.bru`
- [x] 8.7 Create `collection/artpieces/detach_file.bru`
- [x] 8.8 Create `collection/artpieces/set_cover.bru`
- [x] 8.9 Create `collection/files/initiate_upload.bru`
- [x] 8.10 Create `collection/files/complete_upload.bru`

## 9. Lint

- [x] 9.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
