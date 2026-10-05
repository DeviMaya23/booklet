## 1. Backend — Delete standalone images module files

- [x] 1.1 Delete `backend/internal/domain/image.go`
- [x] 1.2 Delete `backend/internal/domain/pending_upload.go`
- [x] 1.3 Delete `backend/internal/handler/image_handler.go` and `image_handler_test.go`
- [x] 1.4 Delete `backend/internal/handler/upload_handler.go` and `upload_handler_test.go`
- [x] 1.5 Delete `backend/internal/repository/image_repository.go` and `image_repository_integration_test.go`
- [x] 1.6 Delete `backend/internal/repository/upload_repository.go` and `upload_repository_integration_test.go`
- [x] 1.7 Delete `backend/internal/usecase/image_repository.go`, `image_usecase.go`, and `image_usecase_test.go`
- [x] 1.8 Delete `backend/internal/usecase/upload_repository.go`, `upload_usecase.go`, and `upload_usecase_test.go`
- [x] 1.9 Delete `backend/internal/worker/generate_thumbnail.go` and `generate_thumbnail_test.go`
- [x] 1.10 Delete `backend/internal/worker/purge_upload.go`

## 2. Backend — Remove image plumbing from character layer

- [x] 2.1 Remove `imageRepo ImageRepository` field and constructor param from `characterUsecase` in `usecase/character_usecase.go`
- [x] 2.2 Remove `GetCharacterImages` method from `characterUsecase`
- [x] 2.3 Remove `ImageRepository` interface from `usecase/character_repository.go` (if present) or wherever it is declared for the character usecase
- [x] 2.4 Remove `GetCharacterImages` handler and `characterImageResponse` type from `handler/character_handler.go`
- [x] 2.5 Remove `GetCharacterImages` from the `CharacterUsecase` interface in `handler/character_handler.go`
- [x] 2.6 Update `usecase/fakes_test.go` — remove `fakeImageRepository` struct and all its methods
- [x] 2.7 Update `usecase/character_usecase_test.go` — remove `&fakeImageRepository{}` arg from all `NewCharacterUsecase` calls

## 3. Backend — Remove image/pending-upload cleanup from user repository

- [x] 3.1 Remove image R2-key collection block from `DeleteAllUserData` in `repository/user_repository.go`
- [x] 3.2 Remove pending-upload R2-key collection block from `DeleteAllUserData`
- [x] 3.3 Remove `DELETE FROM image_characters` and `DELETE FROM images` SQL from `DeleteAllUserData`
- [x] 3.4 Remove `DELETE FROM pending_uploads` from `DeleteAllUserData`

## 4. Backend — Remove wiring from main.go

- [x] 4.1 Remove `imageRepository`, `imageUsecase`, `imageHandler` instantiation
- [x] 4.2 Remove `uploadRepository`, `uploadUsecase`, `uploadHandler` instantiation
- [x] 4.3 Remove `imageRepository` arg from `NewUploadUsecase` call and from `NewCharacterUsecase` call
- [x] 4.4 Remove `GenerateThumbnailWorker` registration (`river.AddWorker`)
- [x] 4.5 Remove `PurgeExpiredUploadsWorker` registration and its periodic job
- [x] 4.6 Remove all `/images` and `/images/:id/*` route registrations (6 routes)
- [x] 4.7 Remove `GET /characters/:id/images` route registration

## 5. Backend — Drop migration

- [x] 5.1 Create `backend/migration/000023_drop_images_module.up.sql` — drop `image_characters`, `pending_uploads`, `images` tables
- [x] 5.2 Create `backend/migration/000023_drop_images_module.down.sql` — recreate the three tables from their original schemas (migrations 003, 004, 005, 018)

## 6. Backend — Quality check

- [x] 6.1 Run `golangci-lint run ./...` from `backend/` and fix any issues

## 7. Frontend — Delete standalone images module files

- [x] 7.1 Delete `frontend/src/features/images/` directory (all API hooks, components)
- [x] 7.2 Delete `frontend/src/pages/ImagesPage.tsx`

## 8. Frontend — Remove images wiring from app shell

- [x] 8.1 Remove `ImagesPage` import and `/app/images` route from `frontend/src/App.tsx`
- [x] 8.2 Remove `{ label: 'Images', to: '/app/images' }` nav item from `frontend/src/components/AppSidebar.tsx`
- [x] 8.3 Replace `/images` fixture URL in `frontend/src/lib/api.test.ts` with `/artpieces`

## 9. Frontend — Quality check

- [x] 9.1 Run `npm run build` from `frontend/` and fix any issues
- [x] 9.2 Run `npm run lint` from `frontend/` and fix any issues

## 10. Collection

- [x] 10.1 Delete `collection/images/` directory
