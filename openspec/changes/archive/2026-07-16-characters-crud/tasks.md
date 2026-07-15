## 1. Migration

- [x] 1.1 Create `backend/migration/000002_create_characters.up.sql` — `characters` table with `id` (UUID PK), `user_id` (FK → users), `name` (text NOT NULL), `hero_image_r2_path` (text nullable), `biography` (text nullable), `is_public` (bool default false), `created_at`, `updated_at`, `deleted_at`
- [x] 1.2 Create `backend/migration/000002_create_characters.down.sql` — drops the `characters` table

## 2. Domain

- [x] 2.1 Create `backend/internal/domain/character.go` — `Character` struct with GORM tags matching the migration schema; include `gorm.DeletedAt` for soft delete

## 3. Repository Interface

- [x] 3.1 Create `backend/internal/usecase/character_repository.go` — export `CharacterRepository` interface with methods: `Create`, `GetByID`, `List`, `Update`, `Delete`
- [x] 3.2 Define `UpdateCharacterParams` struct (pointer fields: `*string` for Name, HeroImageR2Path, Biography; `*bool` for IsPublic) in `character_repository.go` alongside the interface

## 4. Usecase

- [x] 4.1 Create `backend/internal/usecase/character_usecase.go` — `characterUsecase` struct, constructor `NewCharacterUsecase`, and `CreateCharacterParams` struct
- [x] 4.2 Implement `Create` — generate UUID via `uuid.New().String()`, assemble `Character` from params + userID, call repo
- [x] 4.3 Implement `GetByID` — delegate to repo; repo scopes by `user_id`
- [x] 4.4 Implement `List` — delegate to repo scoped by `user_id`
- [x] 4.5 Implement `Update` — delegate `UpdateCharacterParams` to repo
- [x] 4.6 Implement `Delete` — delegate to repo scoped by `user_id`

## 5. Repository Implementation

- [x] 5.1 Create `backend/internal/repository/character_repository.go` — `characterRepository` struct, constructor `NewCharacterRepository`
- [x] 5.2 Implement `Create` — `db.WithContext(ctx).Create(&character)`
- [x] 5.3 Implement `GetByID` — `WHERE id = ? AND user_id = ?`; wrap error with `%w`
- [x] 5.4 Implement `List` — `WHERE user_id = ?` ordered by `created_at DESC`
- [x] 5.5 Implement `Update` — build `map[string]interface{}` from non-nil pointer fields in `UpdateCharacterParams`; call `db.Model(&Character{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)`; wrap error with `%w`
- [x] 5.6 Implement `Delete` — `db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Character{})`; return `gorm.ErrRecordNotFound` if no rows affected

## 6. Handler

- [x] 6.1 Create `backend/internal/handler/character_handler.go` — `CharacterHandler` struct with `CharacterUsecase` interface (5 methods), constructor `NewCharacterHandler`
- [x] 6.2 Implement `CreateCharacter` — bind JSON, validate name non-empty, extract userID via `middleware.AuthenticatedUserIDFromContext`, call usecase, return 201
- [x] 6.3 Implement `GetCharacterByID` — validate UUID path param, extract userID, call usecase, map `gorm.ErrRecordNotFound` → 404
- [x] 6.4 Implement `ListCharacters` — extract userID, call usecase, return 200 with array
- [x] 6.5 Implement `UpdateCharacter` — validate UUID path param, bind JSON to pointer-field struct, extract userID, call usecase, map `gorm.ErrRecordNotFound` → 404, return 200
- [x] 6.6 Implement `DeleteCharacter` — validate UUID path param, extract userID, call usecase, map `gorm.ErrRecordNotFound` → 404, return 204

## 7. Wire-up

- [x] 7.1 In `cmd/server/main.go` `initApp`: instantiate `characterRepository`, `characterUsecase`, `characterHandler`
- [x] 7.2 Register routes on the `protected` group: `POST /characters`, `GET /characters`, `GET /characters/:id`, `PATCH /characters/:id`, `DELETE /characters/:id`

## 8. Bruno Collection

- [x] 8.1 Create `collection/characters/create_character.bru`
- [x] 8.2 Create `collection/characters/get_character.bru`
- [x] 8.3 Create `collection/characters/list_characters.bru`
- [x] 8.4 Create `collection/characters/update_character.bru`
- [x] 8.5 Create `collection/characters/delete_character.bru`

## 9. Usecase Unit Tests

- [x] 9.1 Create `backend/internal/usecase/fakes_test.go` — fake `CharacterRepository` implementing the interface
- [x] 9.2 Create `backend/internal/usecase/character_usecase_test.go`
- [x] 9.3 Write `TestCreate_AssemblesCharacter` — assert UUID is set, UserID matches caller, Name matches input
- [x] 9.4 Write `TestUpdate_PartialFields` — table-driven: name only, is_public false only, multiple fields; assert only the provided fields are passed to the repo (non-nil in params)

## 10. Handler Unit Tests

- [x] 10.1 Create `backend/internal/handler/character_handler_test.go`
- [x] 10.2 Test `CreateCharacter`: happy path (201 + body), missing name (400), malformed JSON (400), usecase error (500)
- [x] 10.3 Test `GetCharacterByID`: happy path (200 + body), not found (404), invalid UUID param (400)
- [x] 10.4 Test `ListCharacters`: happy path (200 + array)
- [x] 10.5 Test `UpdateCharacter`: happy path (200 + body), not found (404), malformed JSON (400)
- [x] 10.6 Test `DeleteCharacter`: happy path (204), not found (404)

## 11. Lint

- [x] 11.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues
