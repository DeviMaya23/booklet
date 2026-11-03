# Tasks

## 1. Data model and repository

- [x] 1.1 Add migration `000037_add_last_used_at_to_characters` (up: add nullable `last_used_at`, backfill `= created_at`; down: drop column). Verify by running migrate up/down locally and checking existing rows have `last_used_at = created_at`.
- [x] 1.2 Add `LastUsedAt *time.Time` to `domain.Character`. Verify `go build ./...` passes.
- [x] 1.3 Add `Sort *string` to `usecase.ListCharacterFilters` and `UpdateLastUsedAt(ctx, ids []uuid.UUID, userID uuid.UUID) error` to the `CharacterRepository` interface; update every implementer and test double (`fakes_test.go` and any others the compiler flags). Verify `go vet ./...` passes.
- [x] 1.4 Implement in `repository/character_repository.go`: `List` orders by `last_used_at DESC NULLS LAST, name ASC` when `Sort` is `recent`, else `created_at DESC`; `UpdateLastUsedAt` bumps all given IDs owned by the user in one query. Add integration tests in `character_repository_integration_test.go` covering: recency order, null-last, name tiebreak, default order unaffected, combination with `q`/`verse_id`, bump of multiple IDs, and no bump for another user's IDs. Verify the integration tests pass.

## 2. Usecases

- [x] 2.1 `character_usecase.Create` sets `LastUsedAt` to now. Add a unit test asserting the created character carries a non-nil `LastUsedAt`. Verify the test passes.
- [x] 2.2 `artpiece_usecase.Create`: after character ownership validation, bump the characters when `CharacterIDs` is non-empty. Add unit tests: bump called with the validated IDs; no bump when empty; no bump when ownership fails (`ErrCharacterNotOwned`). Verify tests pass.
- [x] 2.3 `artpiece_usecase.Update`: add character ownership validation (`GetByIDsAndUserID`, `ErrCharacterNotOwned` on mismatch) then bump, mirroring Create. Add unit tests: bump on non-empty set, no bump on empty set, ownership failure returns `ErrCharacterNotOwned` without bump. Verify tests pass.
- [x] 2.4 `commission_usecase.Create` and `Update`: after the existing character validation, bump the characters when `CharacterIDs` is non-empty. Add unit tests for both: bump on non-empty, no bump on empty, no bump on `ErrCharacterNotOwned`. Verify tests pass.

## 3. Handler and API surface

- [x] 3.1 `character_handler.ListCharacters`: add `c.Validate(&filters)` after `c.Bind`, with `validate:"omitempty,oneof=recent"` on `Sort`; an invalid value returns 422 via the project's validation error shape. Add handler unit tests: `sort=recent` is passed through to the usecase filters, absent `sort` passes nil, `sort=bogus` returns 422 with the structured error. Verify tests pass.
- [x] 3.2 Update `collection/characters/list_characters.bru` with an optional `~sort: recent` query param. Verify the file parses (Bruno loads the collection) and the request returns 200 with and without the param.
- [x] 3.3 Add `last_used_at` to the `characters` table in `DATABASE.md`. Verify the mermaid block still renders.
- [x] 3.4 Run `golangci-lint run` in `backend/` and fix every reported issue. Verify it exits clean.

## 4. Frontend

- [x] 4.1 `useCharacters(options?: { sort?: 'recent' })`: append `?sort=recent` when set and use query key `['characters', 'recent']` for that variant; default call and `CHARACTERS_QUERY_KEY` unchanged. Add tests at the API boundary (mock the fetch wrapper, fresh `QueryClientProvider` per test): default call hits `/characters` with no param; the recent variant hits `/characters?sort=recent`; invalidating `CHARACTERS_QUERY_KEY` refetches both variants. Verify the tests pass.
- [x] 4.2 `ArtpieceFormModal`, `ArtpieceDetailView`, `ArtistCharacterFilter`: call `useCharacters({ sort: 'recent' })`. Add or extend component tests asserting suggestions render in the order the API returns (A, B, C) with selected characters excluded, and that the request includes `sort=recent`. Verify the tests pass.
- [x] 4.3 `useCreateArtpiece`, `useUpdateArtpiece`, `useCreateCommission`, `useUpdateCommission`: invalidate `CHARACTERS_QUERY_KEY` on success. Extend each hook's test to assert the characters query is invalidated. Verify the tests pass.
- [x] 4.4 Confirm `CharactersPage`, `ChangeThumbnailModal`, and `useArtpiecesFilter` still call `useCharacters()` without `sort` (grep every `useCharacters(` call site). Verify by listing the call sites in the PR description.
- [x] 4.5 Run `npm run build` and `npm run lint` in `frontend/` and fix every reported issue. Verify both exit clean.

## 5. Integration check

- [x] 5.1 Manual run against a local stack: create two artpieces using different characters, open the artpiece create modal and the gallery filter, and confirm the most recently used character is listed first in all three pickers while the Characters page order is unchanged. Verify by observation.
