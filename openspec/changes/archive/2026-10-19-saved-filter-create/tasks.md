# Tasks

## 1. BE: Migration

- [x] 1.1 Create `backend/migration/000032_saved_filters_unique_name.up.sql` adding a `UNIQUE (user_id, name)` constraint on `saved_filters`, and the corresponding `.down.sql` that drops it; verify the migration applies cleanly with `golang-migrate` or the project's migration runner

## 2. BE: Usecase layer

- [x] 2.1 Add `ErrDuplicateSavedFilterName` sentinel error in `usecase/saved_filter_repository.go` alongside `ErrSavedFilterNotFound`
- [x] 2.2 Add a narrow `ArtpieceThumbnailRepository` interface to the saved filter usecase (one method: `GetByID(ctx, id, userID) (*domain.Artpiece, error)`), placed in `usecase/saved_filter_repository.go` per interface file-placement rules; verify it is satisfied by the existing artpiece repository implementation without changes
- [x] 2.3 Update `CreateSavedFilterParams` to replace `ThumbnailR2Path *string` with `ThumbnailArtpieceID *uuid.UUID`; update `savedFilterUsecase` constructor to accept `ArtpieceThumbnailRepository`; update `Create` to look up the artpiece when `ThumbnailArtpieceID` is set, extract `CoverFile.ThumbnailR2Path`, and pass it to the repo — if the artpiece has no cover thumbnail, store `nil`; map `ErrDuplicateSavedFilterName` from the repo to the sentinel; verify `go build ./...` passes
- [x] 2.4 Write usecase unit tests in `saved_filter_usecase_test.go`: (a) `Create` with a valid artpiece ID copies the R2 path from the artpiece into the saved filter params; (b) `Create` when the repo returns the duplicate-name error maps it to `ErrDuplicateSavedFilterName`; verify `go test ./internal/usecase/... -run TestSavedFilter` passes

## 3. BE: Repository layer

- [x] 3.1 Update `saved_filter_repository.go` `Create` method to detect a unique constraint violation (using the `isUniqueConstraintViolation` helper pattern from `artist_repository.go`) and return `ErrDuplicateSavedFilterName`; verify `go build ./...` passes
- [x] 3.2 Add an integration test in `saved_filter_repository_integration_test.go` for the duplicate-name scenario: create two saved filters with the same name for the same user and assert `ErrDuplicateSavedFilterName` is returned; also assert two different users can each have a filter with the same name; verify the integration test passes

## 4. BE: Handler layer

- [x] 4.1 Add `Presigner` to `SavedFilterHandler` — extend `NewSavedFilterHandler` to accept a `Presigner` parameter; update `toSavedFilterResponse` to sign `ThumbnailR2Path` and return `ThumbnailURL *string` in place of `ThumbnailR2Path`; update `savedFilterResponse` struct accordingly; verify `go build ./...` passes
- [x] 4.2 Update `createSavedFilterRequest` to replace `thumbnail_r2_path` with `thumbnail_artpiece_id` (string, optional, no `validate` tag); update `CreateSavedFilter` handler to parse the artpiece UUID, pass `ThumbnailArtpieceID` to the usecase, and map `ErrDuplicateSavedFilterName` to a 422 with a `name` field-level error matching the `validationErrResponse` shape; verify `go build ./...` passes
- [x] 4.3 Update `main.go` to pass the presigner and artpiece repository to `NewSavedFilterUsecase` and `NewSavedFilterHandler`; verify `go build ./...` passes
- [x] 4.4 Update `saved_filter_handler_test.go`: add the presigner spy (mirror the pattern used in `artpiece_handler_test.go`); add test for duplicate name → 422 with `name` field error; add test that the `CreateSavedFilter` response contains `thumbnail_url` and not `thumbnail_r2_path`; update any existing tests whose assertions reference `thumbnail_r2_path`; verify `go test ./internal/handler/... -run TestSavedFilter` passes

## 5. BE: Bruno

- [x] 5.1 Update `collection/saved_filters/create_saved_filter.bru` to replace `thumbnail_r2_path` with `thumbnail_artpiece_id` in the request body and add a comment showing the expected `thumbnail_url` field in the response; verify the file is valid Bruno syntax

## 6. FE: `ArtistCharacterFilter` — footerAction slot

- [x] 6.1 Add optional `footerAction?: ReactNode` prop to `ArtistCharacterFilter`; render it in the popover footer row to the right of the "Clear all" link (visible only when the footer row is rendered); verify callers `CommissionFormModal` and `AddToExistingArtpieceModal` compile with no changes
- [x] 6.2 Add two test scenarios to `ArtistCharacterFilter.test.tsx`: (a) when `footerAction` is provided the slot content is rendered; (b) when `footerAction` is omitted the slot content is absent; verify `npm test -- ArtistCharacterFilter` passes

## 7. FE: Filter URL params

- [x] 7.1 Refactor `useArtpiecesFilter` to sync filter state with `useSearchParams`; params: `artist` (UUID string), `characters` (comma-separated UUID strings), `match` (`all` | `any`); on mount, hydrate artist and character objects from the cached `useArtists` / `useCharacters` query data; clearing all filters removes params from the URL; verify `ArtpiecesPage` still compiles and filter interactions update the URL in the browser
- [x] 7.2 Update `useArtpiecesFilter.test.ts` to cover: (a) URL params present on mount → filter state initialised from params; (b) changing a filter → URL updates; (c) clearing filters → params removed; verify `npm test -- useArtpiecesFilter` passes

## 8. FE: API hook — create saved filter

- [x] 8.1 Create `frontend/src/features/artpieces/api/useCreateSavedFilter.ts` — a React Query mutation hook that POSTs to `/saved_filters` with `{ name, thumbnail_artpiece_id, filter_payload }`; return type should include the created saved filter (with `thumbnail_url`); verify `npm run build` passes with no type errors

## 9. FE: `SaveFilterModal`

- [x] 9.1 Create `frontend/src/features/artpieces/components/SaveFilterModal.tsx` implementing all `web-save-filter-modal` spec requirements: filter summary chip, name input, thumbnail picker grid (up to 10 artpieces, first pre-selected), no-results empty state with disabled Save, dashboard preview tile, Save/Cancel buttons; use `useCreateSavedFilter` for the mutation; call `onSuccess(savedFilter)` on 201; show inline name error on 422; show an error toast on other server errors; verify `npm run build` passes
- [x] 9.2 Write `SaveFilterModal.test.tsx` colocated with the component; test: (a) no-results state shows empty state and Save is disabled; (b) submitting a duplicate name shows the inline field error and keeps modal open; (c) Cancel closes the modal without calling the mutation; (d) successful save calls `onSuccess` and shows a success toast; verify `npm test -- SaveFilterModal` passes

## 10. FE: `ArtpiecesPage` — wiring

- [x] 10.1 In `ArtpiecesPage`, manage the Save Filter modal's open state and ephemeral "saved as" snapshot (`{ name, artistId, characterIds } | null`); pass a `footerAction` to `ArtistCharacterFilter` that renders either the "Save to dashboard" bookmark button (when filters are active and the current filter doesn't match the snapshot) or the "Saved as '<name>'" confirmation text (when the snapshot matches); reset the snapshot when filter values change; open `SaveFilterModal` on button click; on modal `onSuccess`, store the snapshot and show a success toast; verify the button and confirmation appear and disappear correctly in the browser

## 11. FE: Build and lint

- [x] 11.1 Run `npm run build` from `frontend/` and fix all type errors and build warnings
- [x] 11.2 Run `npm run lint` from `frontend/` and fix all lint findings without using `eslint-disable` suppressions

## 12. BE: Lint

- [x] 12.1 Run `golangci-lint run ./...` from `backend/` and fix all reported findings

## Workflow follow-up

- Archive the change once the PR is merged and reviewed.
