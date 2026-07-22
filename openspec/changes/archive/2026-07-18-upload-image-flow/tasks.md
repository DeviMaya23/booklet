## 1. Database Migrations

- [x] 1.1 Create `backend/migration/000004_create_pending_upload.up.sql` — `pending_upload` table with columns: `id uuid PK`, `user_id text NOT NULL FK users(id)`, `r2_key text NOT NULL`, `mime_type text NOT NULL`, `title text`, `artist_name text`, `artist_link text`, `notes text`, `character_ids jsonb NOT NULL DEFAULT '[]'`, `created_at timestamptz NOT NULL DEFAULT NOW()`; create index on `user_id`
- [x] 1.2 Create `backend/migration/000004_create_pending_upload.down.sql` — drop `pending_upload` table
- [x] 1.3 Create `backend/migration/000005_add_mime_type_to_images.up.sql` — `ALTER TABLE images ADD COLUMN mime_type TEXT NOT NULL DEFAULT 'image/jpeg'`; then `ALTER TABLE images ALTER COLUMN mime_type DROP DEFAULT`
- [x] 1.4 Create `backend/migration/000005_add_mime_type_to_images.down.sql` — `ALTER TABLE images DROP COLUMN mime_type`

## 2. Domain

- [x] 2.1 Add `MimeType string` field to `domain.Image` with gorm tag `type:text;not null;column:mime_type`
- [x] 2.2 Create `internal/domain/pending_upload.go` — `PendingUpload` struct with fields matching the migration schema; `CharacterIDs []string` with gorm tag `type:jsonb;serializer:json;column:character_ids`

## 3. Health Check — R2 Probe

- [x] 3.1 In `internal/handler/health.go`: add `R2Pinger` interface (`Ping(ctx context.Context) error`); add it as a field on `HealthHandler`; uncomment and wire the r2 probe in `NewHealthHandler` and `GetHealth`; add `R2` field to `healthResponse`
- [x] 3.2 In `cmd/server/main.go` `initApp`: instantiate `storage.NewR2Storage(cfg.R2, tel)` and pass it to `NewHealthHandler`
- [x] 3.3 Write unit tests for `HealthHandler` covering: all healthy, DB probe fails, R2 probe fails

## 4. Upload Usecase Interfaces

- [x] 4.1 Create `internal/usecase/upload_repository.go` — define `StorageService` interface (`GeneratePresignedPutURL(ctx, key, contentType string, ttl time.Duration) (string, error)`); define `UploadRepository` interface (`CreatePendingUpload(ctx context.Context, p *domain.PendingUpload) error` and `CompleteUpload(ctx context.Context, pendingID, userID string, validCharIDs []uuid.UUID) (*domain.Image, error)`); define `FindCharactersByIDsAndUserParams` and any shared param types needed

## 5. Upload Usecase

- [x] 5.1 Create `internal/usecase/upload_usecase.go` — implement `InitialUpload` (generate UUID, derive ext from mime_type, build R2 key, call `StorageService.GeneratePresignedPutURL`, persist `PendingUpload`, return id + url + expires_at) and `CompleteUpload` (fetch pending upload, query valid character IDs scoped to user via `CharacterRepository`, call `UploadRepository.CompleteUpload` with filtered IDs); add `// TODO: trigger thumbnail generation job` stub in CompleteUpload after the repo call
- [x] 5.2 Write unit tests for `upload_usecase.go` — CompleteUpload character filtering: some IDs valid → only valid ones passed to repo; all IDs invalid → empty slice passed to repo; InitialUpload R2 key format: assert the generated key matches `users/{userID}/images/{uuid}.{ext}` pattern

## 6. Upload Repository Implementation

- [x] 6.1 Create `internal/repository/upload_repository.go` — implement `CreatePendingUpload` (insert row); implement `CompleteUpload` as a single GORM transaction: fetch and verify pending_upload ownership (return `gorm.ErrRecordNotFound` if missing or wrong user), delete pending_upload row, insert image row (copying `r2_key` → `image_r2_path`, `mime_type`, metadata fields), associate characters via `image_characters`, return the created image with Characters preloaded
- [x] 6.2 Write integration tests for `upload_repository.go` — `CreatePendingUpload`: row persisted with correct fields; `CompleteUpload` success: pending deleted, image created with correct fields and character associations; `CompleteUpload` wrong user: returns `gorm.ErrRecordNotFound`; `CompleteUpload` not found: returns `gorm.ErrRecordNotFound`

## 7. Image Handler and Repository Updates

- [x] 7.1 Add `MimeType string` to `imageResponse` struct in `internal/handler/image_handler.go` and map it in `toImageResponse`
- [x] 7.2 Update `internal/repository/image_repository_integration_test.go` seed data to include `mime_type` on all inserted image rows

## 8. Upload Handler

- [x] 8.1 Create `internal/handler/upload_handler.go` — define `UploadUsecase` interface; implement `InitialUpload` handler: bind + validate request using Echo validator (`mime_type` uses `validate:"required,oneof=image/jpeg image/png"`), extract userID from context, call usecase, return 201 with `{ id, upload_url, expires_at }`; implement `CompleteUpload` handler: validate UUID path param, call usecase, return 201 with image response
- [x] 8.2 Write unit tests for `upload_handler.go` — `InitialUpload`: happy path (assert 201 + response fields), missing mime_type → 422, invalid mime_type → 422, malformed JSON → 400; `CompleteUpload`: happy path (assert 201 + image body), pending not found → 404, invalid UUID path param → 400

## 9. Wire Up Routes

- [x] 9.1 In `cmd/server/main.go` `initApp`: instantiate `uploadRepository`, `uploadUsecase` (inject `StorageService` and `CharacterRepository`), `uploadHandler`; register `protected.POST("/images", uploadHandler.InitialUpload)` and `protected.POST("/images/:id/complete", uploadHandler.CompleteUpload)`

## 10. Bruno Collection

- [x] 10.1 Create `collection/images/initial_upload.bru` — `POST /images` with JSON body containing `mime_type`, `title`, `artist_name`, `artist_link`, `notes`, `character_ids`
- [x] 10.2 Create `collection/images/complete_upload.bru` — `POST /images/:id/complete` with path param `id`

## 11. Lint

- [x] 11.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues
