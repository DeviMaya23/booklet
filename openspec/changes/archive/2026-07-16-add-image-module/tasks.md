## 1. Database Migration

- [x] 1.1 Create migration: `images` table (`id uuid PK`, `user_id text NOT NULL`, `image_r2_path text NOT NULL`, `thumbnail_r2_path text`, `artist_name text`, `artist_link text`, `notes text`, `created_at timestamptz NOT NULL`, `updated_at timestamptz NOT NULL`)
- [x] 1.2 Create migration: `image_characters` join table (`image_id uuid FK → images`, `character_id uuid FK → characters`, composite PK)

## 2. Domain

- [x] 2.1 Create `internal/domain/image.go` — `Image` struct with GORM tags, including `Characters []Character` field tagged `gorm:"many2many:image_characters;"`

## 3. Usecase Layer

- [x] 3.1 Create `internal/usecase/image_repository.go` — `ImageRepository` interface and `UpdateImageParams` struct
- [x] 3.2 Create `internal/usecase/image_usecase.go` — `imageUsecase` with `GetByID`, `List`, `Update`, `Delete`; each method starts a telemetry span and sets span error on failure
- [x] 3.3 Write `internal/usecase/image_usecase_test.go` — unit tests per CONVENTIONS.md (only scenarios with real logic worth testing)

## 4. Repository Layer

- [x] 4.1 Create `internal/repository/image_repository.go` — GORM implementation of `ImageRepository`; use `Preload("Characters")` on all reads; implement `GetByID`, `List`, `Update` (with character ownership validation + `Association("Characters").Replace`), `Delete` (hard delete)

## 5. Handler Layer

- [x] 5.1 Create `internal/handler/image_handler.go` — `ImageHandler` with `GetImageByID`, `ListImages`, `UpdateImage`, `DeleteImage`; each method starts a telemetry span; responses use snake_case via `imageResponse` struct
- [x] 5.2 Write `internal/handler/image_handler_test.go` — unit tests covering happy path, each error mapping, UUID validation, and malformed JSON per CONVENTIONS.md

## 6. Wiring

- [x] 6.1 Update `cmd/server/main.go` `initApp` — instantiate `imageRepository`, `imageUsecase`, `imageHandler` and register routes: `GET /images`, `GET /images/:id`, `PATCH /images/:id`, `DELETE /images/:id` on the protected group

## 7. Bruno Collection

- [x] 7.1 Create `collection/images/` directory with Bruno request files: `list_images.bru`, `get_image.bru`, `update_image.bru`, `delete_image.bru`

## 8. Lint

- [x] 8.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
