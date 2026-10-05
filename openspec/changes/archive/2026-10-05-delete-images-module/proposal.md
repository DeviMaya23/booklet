## Why

The images module (upload, list, get, update, delete) is being phased out in favour of the artpieces + files model. The app is not live, so there are no real users or data to migrate — this is a clean removal.

## What Changes

- **BREAKING** Remove all `/images` and `/images/:id/*` API endpoints (GET, PUT, DELETE, POST init/complete)
- **BREAKING** Remove `GET /characters/:id/images` endpoint (will be re-added in a future change against artpieces/files)
- Delete `images`, `image_characters`, and `pending_uploads` DB tables via a drop migration
- Delete all BE layers for images: domain model, repository, usecase, handler, workers (`generate_thumbnail`, `purge_expired_uploads`)
- Delete all FE layers: `features/images/` directory, `ImagesPage`, route `/app/images`, sidebar nav item
- Remove image/pending-upload cleanup from `DeleteAllUserData` in the user repository
- Remove all related Bruno collection files under `collection/images/`

## Capabilities

### New Capabilities

_None — this change only removes capabilities._

### Modified Capabilities

- `image-management`: being deleted entirely
- `image-upload`: being deleted entirely
- `image-thumbnail-generation`: being deleted entirely
- `stale-upload-purge`: purge of `pending_uploads` table being removed; file-upload purge remains
- `character-images`: `GET /characters/:id/images` endpoint removed; character-to-image relationship dropped
- `account-deletion`: image and pending-upload R2 key collection + DB row deletion removed from user teardown
- `web-images`: FE images page and nav item removed
- `web-images-create`: FE image upload modal removed
- `web-images-edit`: FE image edit modal removed

## Impact

**Backend**
- `backend/internal/domain/image.go`, `pending_upload.go`
- `backend/internal/handler/image_handler.go`, `upload_handler.go` (+ tests)
- `backend/internal/repository/image_repository.go`, `upload_repository.go` (+ integration tests)
- `backend/internal/usecase/image_repository.go`, `image_usecase.go`, `upload_repository.go`, `upload_usecase.go` (+ tests)
- `backend/internal/worker/generate_thumbnail.go`, `purge_upload.go` (+ tests)
- `backend/cmd/server/main.go` — remove handler wiring, route registrations, worker/periodic-job registrations
- `backend/internal/handler/character_handler.go` — remove `GetCharacterImages`
- `backend/internal/usecase/character_usecase.go` — remove `imageRepo` field + `GetCharacterImages`
- `backend/internal/repository/user_repository.go` — remove image/pending-upload sections from `DeleteAllUserData`
- `backend/internal/usecase/fakes_test.go` — remove `fakeImageRepository`
- `backend/migration/` — new `000023_drop_images_module` migration pair

**Frontend**
- `frontend/src/features/images/` (entire directory)
- `frontend/src/pages/ImagesPage.tsx`
- `frontend/src/App.tsx` — remove import + route
- `frontend/src/components/AppSidebar.tsx` — remove nav item
- `frontend/src/lib/api.test.ts` — replace `/images` fixture URL with another path

**Collection**
- `collection/images/` (entire directory)
