# Tasks

## 1. Migrations

- [x] 1.1 Create `backend/migration/000033_create_verses.up.sql` with the `verses` table (`id`, `user_id`, `name`, `last_used_at`, `created_at`, `updated_at`, `UNIQUE(user_id, name)`, index on `user_id`) and `000033_create_verses.down.sql` that drops the table; verify the migration applies and rolls back cleanly with `make migrate-up` / `make migrate-down`
- [x] 1.2 Create `backend/migration/000034_add_verse_id_to_characters.up.sql` adding `verse_id uuid NULL REFERENCES verses(id) ON DELETE SET NULL` to `characters` and `000034_add_verse_id_to_characters.down.sql` that drops the column; verify apply and rollback

## 2. Domain models

- [x] 2.1 Add `domain.Verse` struct to `backend/internal/domain/` with fields: `ID uuid.UUID`, `UserID uuid.UUID`, `Name string`, `LastUsedAt *time.Time`, `CreatedAt time.Time`, `UpdatedAt time.Time`; verify it compiles
- [x] 2.2 Add `VerseID *uuid.UUID` column field to `domain.Character` with GORM tag `gorm:"type:uuid;column:verse_id"`; verify it compiles

## 3. Verse repository layer

- [x] 3.1 Create `backend/internal/usecase/verse_repository.go` defining `VerseRepository` interface (`Create`, `GetByID`, `List`, `Update`, `Delete`, `UpdateLastUsedAt`), `CreateVerseParams`, `UpdateVerseParams`, `ListVerseFilters`, and `ErrVerseNameConflict` / `ErrVerseNotOwned` sentinel errors; verify it compiles
- [x] 3.2 Create `backend/internal/repository/verse_repository.go` implementing `VerseRepository`; `List` orders by `last_used_at DESC NULLS LAST, name ASC`; `Create` / `Update` detect unique-constraint violations (pgx error code `23505`) and return `ErrVerseNameConflict`; `GetByID` / `Delete` scope to `user_id` and return `gorm.ErrRecordNotFound` when no row matches; verify it compiles
- [x] 3.3 Create `backend/internal/repository/verse_repository_integration_test.go` with integration tests for: Create (success, duplicate name returns `ErrVerseNameConflict`), GetByID (own verse, another user's verse returns not-found), List (ordering with mixed `last_used_at` values), Update (success, rename to existing name returns `ErrVerseNameConflict`, not found), Delete (success returns 204-equivalent, sets `verse_id` to null on any character that referenced it), UpdateLastUsedAt; verify tests pass

## 4. Verse usecase and handler

- [x] 4.1 Create `backend/internal/usecase/verse_usecase.go` implementing Create, GetByID, List, Update, Delete; verify it compiles
- [x] 4.2 Create `backend/internal/usecase/verse_usecase_test.go` with unit tests covering: Create success, Create with name conflict returns `ErrVerseNameConflict`, Update not found, Update name conflict, Delete not found; verify tests pass
- [x] 4.3 Create `backend/internal/handler/verse_handler.go` with `VersHandler` for `POST /verses`, `GET /verses`, `GET /verses/:id`, `PATCH /verses/:id`, `DELETE /verses/:id`; use `c.Bind` + `c.Validate` for body params; `PATCH` body requires `name` (non-empty); responses match the verse response shape from the spec; `ErrVerseNameConflict` → 409, `gorm.ErrRecordNotFound` → 404; verify it compiles
- [x] 4.4 Create `backend/internal/handler/verse_handler_test.go` with unit tests covering: POST success (201), POST with missing name (422), POST with duplicate name (409), GET list (200), GET by ID (200), GET by ID not found (404), PATCH success (200), PATCH empty name (422), PATCH name conflict (409), PATCH not found (404), DELETE success (204), DELETE not found (404); verify tests pass
- [x] 4.5 Create Bruno request files under `collection/verses/`: `create_verse.bru`, `list_verses.bru`, `get_verse.bru`, `update_verse.bru`, `delete_verse.bru`; verify they open in Bruno without errors

## 5. Character layer updates

- [x] 5.1 Add `VerseID *uuid.UUID` to `usecase.CreateCharacterParams` and `usecase.UpdateCharacterParams` in `backend/internal/usecase/character_repository.go`; verify it compiles
- [x] 5.2 Update `characterUsecase` in `backend/internal/usecase/character_usecase.go`: add `verseRepo VerseRepository` field and constructor param; in `Create` and `Update`, when `params.VerseID` is non-nil call `verseRepo.GetByID` scoped to `userID` (return `ErrVerseNotOwned` on miss), then proceed with character persist, then call `verseRepo.UpdateLastUsedAt`; verify it compiles
- [x] 5.3 Update `backend/internal/usecase/character_usecase_test.go` with unit tests for: Create with valid `verse_id` (verse validated + `last_used_at` bumped), Create with unowned `verse_id` (returns `ErrVerseNotOwned` before persisting character), Update with valid `verse_id`, Update with unowned `verse_id`; verify tests pass
- [x] 5.4 Add `VerseID *uuid.UUID` to `usecase.ListCharacterFilters` with `query:"verse_id"` tag; update `characterRepository.List` in `backend/internal/repository/character_repository.go` to apply `WHERE verse_id = ?` when the filter is non-nil; verify it compiles
- [x] 5.5 Add an integration test to `backend/internal/repository/character_repository_integration_test.go` for `List` with a `verse_id` filter: seed two characters, one with a verse and one without, assert only the verse-assigned character is returned; verify test passes
- [x] 5.6 Update `backend/internal/handler/character_handler.go`: add `VerseID *uuid.UUID` to `createCharacterRequest` and `updateCharacterRequest`; add `verse_id *string` to `ListCharacterFilters` binding and parse with `uuid.Parse` before passing to usecase; pass `VerseID` through to usecase params; add `VerseID *string` to `characterResponse` (UUID string or null); update `toCharacterResponse` to populate it from `character.VerseID`; map `ErrVerseNotOwned` → 422 in Create and Update handlers; verify it compiles
- [x] 5.7 Update `backend/internal/handler/character_handler_test.go` with unit tests for: Create with valid `verse_id` (201), Create with unowned `verse_id` (422), `GET /characters?verse_id=<uuid>` binding (200), Update with valid `verse_id` (200), Update with unowned `verse_id` (422); verify tests pass

## 6. Wire up in main.go

- [x] 6.1 In `backend/cmd/server/main.go` add `verseRepository := repository.NewVerseRepository(db)`, `verseUsecase := usecase.NewVerseUsecase(verseRepository, tel)`, `verseHandler := httphandler.NewVerseHandler(verseUsecase, tel)`; pass `verseRepository` to `usecase.NewCharacterUsecase`; register routes `POST /verses`, `GET /verses`, `GET /verses/:id`, `PATCH /verses/:id`, `DELETE /verses/:id`; verify the server builds and starts cleanly

## 7. Lint

- [x] 7.1 Run `golangci-lint run ./...` from `backend/` and fix all reported issues; verify it exits 0

## Workflow follow-up

- Archive the change after review.
