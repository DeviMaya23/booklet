## 1. Migration

- [x] 1.1 Create `000028_create_saved_filters.up.sql` — `CREATE TABLE saved_filters` with columns: `id uuid PK`, `user_id uuid NOT NULL FK -> users`, `name text NOT NULL`, `thumbnail_r2_path text NULL`, `filter_payload jsonb NOT NULL`, `created_at timestamptz NOT NULL DEFAULT now()`, `updated_at timestamptz NOT NULL DEFAULT now()`
- [x] 1.2 Create `000028_create_saved_filters.down.sql` — `DROP TABLE IF EXISTS saved_filters`

## 2. Domain

- [x] 2.1 Add `domain.SavedFilter` struct to `backend/internal/domain/` — fields: `ID uuid`, `UserID uuid`, `Name string`, `ThumbnailR2Path *string`, `FilterPayload json.RawMessage` (`gorm:"type:jsonb"`), `CreatedAt time.Time`, `UpdatedAt time.Time`

## 3. Repository

- [x] 3.1 Define `SavedFilterRepository` interface in `backend/internal/usecase/` with methods: `Create`, `GetByID`, `List`, `Update`, `Delete`
- [x] 3.2 Implement `savedFilterRepository` in `backend/internal/repository/saved_filter_repository.go` — `Create`, `GetByID` (scoped by userID), `List` (scoped by userID, ORDER BY created_at DESC), `Update` (PATCH semantics — only set fields), `Delete` (scoped by userID, hard delete)

## 4. Usecase

- [x] 4.1 Define param types: `CreateSavedFilterParams`, `UpdateSavedFilterParams` (with `Patch[string]` for name and thumbnail, `*json.RawMessage` for payload)
- [x] 4.2 Implement `savedFilterUsecase` in `backend/internal/usecase/saved_filter_usecase.go` — `Create`, `GetByID`, `List`, `Update`, `Delete`; wrap `gorm.ErrRecordNotFound` as `ErrSavedFilterNotFound`
- [x] 4.3 Write usecase unit tests in `backend/internal/usecase/saved_filter_usecase_test.go` using an in-memory fake repository — test scenarios: `Create` assembles the record with correct userID; `GetByID` returns `ErrSavedFilterNotFound` when repo returns not-found; `Update` returns `ErrSavedFilterNotFound` when repo returns not-found; `Delete` returns `ErrSavedFilterNotFound` when repo returns not-found

## 5. Handler

- [x] 5.1 Implement `SavedFilterHandler` in `backend/internal/handler/saved_filter_handler.go` — `CreateSavedFilter`, `ListSavedFilters`, `GetSavedFilter`, `PatchSavedFilter`, `DeleteSavedFilter`; request structs use `c.Bind` + `c.Validate`; PATCH request uses `Patch[string]` for name/thumbnail and `*json.RawMessage` for payload; response struct uses `json.RawMessage` for filter_payload (inline embed, no double-encode)
- [x] 5.2 Write handler unit tests in `backend/internal/handler/saved_filter_handler_test.go` using a spy usecase — test scenarios: `CreateSavedFilter` returns 422 when name is missing; `CreateSavedFilter` happy path returns 201 with body; `GetSavedFilter` returns 404 when usecase returns `ErrSavedFilterNotFound`; `PatchSavedFilter` returns 404 when usecase returns `ErrSavedFilterNotFound`; `DeleteSavedFilter` returns 404 when usecase returns `ErrSavedFilterNotFound`; `DeleteSavedFilter` happy path returns 204

## 6. Route Registration

- [x] 6.1 Wire up `SavedFilterHandler` in `backend/cmd/server/main.go` — instantiate repository, usecase, handler; register routes on the authenticated group: `POST /saved_filters`, `GET /saved_filters`, `GET /saved_filters/:id`, `PATCH /saved_filters/:id`, `DELETE /saved_filters/:id`

## 7. Bruno Collection

- [x] 7.1 Create `collection/saved_filters/` directory with `.bru` files for all five endpoints: `create_saved_filter.bru`, `list_saved_filters.bru`, `get_saved_filter.bru`, `patch_saved_filter.bru`, `delete_saved_filter.bru`

## 8. Lint

- [x] 8.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
