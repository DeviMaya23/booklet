## 1. Migration

- [x] 1.1 Create `000010_add_artists_entity.up.sql`: create `artists` table with `UNIQUE(user_id, name)` constraint
- [x] 1.2 Add `artist_id UUID REFERENCES artists(id) ON DELETE SET NULL` to `images`; drop `artist_name`, `artist_link` columns
- [x] 1.3 Add `artist_id UUID REFERENCES artists(id) ON DELETE SET NULL` to `pending_uploads`; drop `artist_name`, `artist_link` columns
- [x] 1.4 Create `000010_add_artists_entity.down.sql`: reverse the migration (drop `artists`, remove `artist_id`, re-add `artist_name`/`artist_link`)

## 2. Domain

- [x] 2.1 Create `backend/internal/domain/artist.go` with `Artist` struct (GORM tags, no soft delete)
- [x] 2.2 Update `backend/internal/domain/image.go`: replace `ArtistName`/`ArtistLink` fields with `ArtistID *uuid.UUID` and `Artist *Artist`
- [x] 2.3 Update `backend/internal/domain/pending_upload.go`: replace `ArtistName`/`ArtistLink` with `ArtistID *uuid.UUID`

## 3. Artist Repository & Usecase

- [x] 3.1 Define `ArtistRepository` interface in `backend/internal/usecase/artist_repository.go` (Create, GetByID, GetByIDAndUserID, List, Update, Delete)
- [x] 3.2 Create `backend/internal/repository/artist_repository.go` implementing the interface; handle unique constraint violation as a named error
- [x] 3.3 Create `backend/internal/usecase/artist_usecase.go` with CRUD methods; wrap `gorm.ErrRecordNotFound` on not found (passed through from repo), `ErrArtistNameConflict` on duplicate name
- [x] 3.4 Write unit tests for `artist_usecase.go` covering: create success, create duplicate name conflict, get not found, update name conflict, delete success

## 4. Artist Handler

- [x] 4.1 Create `backend/internal/handler/artist_handler.go` with `CreateArtist`, `ListArtists`, `GetArtistByID`, `UpdateArtist`, `DeleteArtist`
- [x] 4.2 Map `ErrArtistNameConflict` → 409; `ErrArtistNotOwned` / not found → 404; validation errors → 422
- [x] 4.3 Write handler tests for `artist_handler.go` covering: create success + conflict, list, get + not found, update + conflict, delete + not found

## 5. Image Repository & Usecase Updates

- [x] 5.1 Update `UpdateImageParams` in `backend/internal/usecase/image_repository.go`: replace `ArtistName`/`ArtistLink` with `ArtistID **uuid.UUID` (pointer-to-pointer for PATCH null semantics)
- [x] 5.2 Update `backend/internal/repository/image_repository.go`: `Update` handles `artist_id` field; `GetByID` and `List` add `.Preload("Artist")`; validate artist ownership when `artist_id` is non-nil (return `ErrArtistNotOwned` on failure)
- [x] 5.3 Update `backend/internal/usecase/image_usecase.go`: pass `ArtistID` through `Update`

## 6. Upload Usecase Update

- [x] 6.1 Update `InitialUploadParams` in `backend/internal/usecase/upload_usecase.go`: replace `ArtistName`/`ArtistLink` with `ArtistID *uuid.UUID`
- [x] 6.2 Update `CompleteUpload` in `upload_usecase.go`: if `pending.ArtistID` is non-nil, query artist by `(id, user_id)`; set to nil if not found; pass resolved `ArtistID` when creating the `Image`

## 7. Handler Updates

- [x] 7.1 Update `backend/internal/handler/upload_handler.go`: replace `ArtistName`/`ArtistLink` in `initialUploadRequest` with `ArtistID *string` (UUID shape validation); parse and pass `ArtistID` to usecase
- [x] 7.2 Update `backend/internal/handler/image_handler.go`: replace `ArtistName`/`ArtistLink` in `updateImageRequest` with `ArtistID json.RawMessage`; update `imageResponse` to include `artist_id` + `artist_name`; update `toImageResponse` to read from `image.Artist`
- [x] 7.3 Update existing handler tests in `upload_handler_test.go` and `image_handler_test.go` to reflect new request/response fields

## 8. Routing

- [x] 8.1 Wire `ArtistRepository`, `ArtistUsecase`, `ArtistHandler` in `backend/cmd/server/main.go`
- [x] 8.2 Register artist routes under the protected group: `POST /artists`, `GET /artists`, `GET /artists/:id`, `PATCH /artists/:id`, `DELETE /artists/:id`

## 9. Bruno

- [x] 9.1 Create Bruno request files for all 5 artist endpoints under the appropriate Bruno collection folder

## 10. Lint & Build

- [x] 10.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
