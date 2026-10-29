# Tasks

## 1. BE — Migration: add sort_order to verses

- [x] 1.1 Create `backend/migration/000035_add_sort_order_to_verses.up.sql`: `ALTER TABLE verses ADD COLUMN sort_order INT NOT NULL DEFAULT 0`, then backfill with `ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at) - 1`. Verify by inspecting the migration SQL for correctness.
- [x] 1.2 Create `backend/migration/000035_add_sort_order_to_verses.down.sql`: `ALTER TABLE verses DROP COLUMN sort_order`. Verify the down migration reverts the up migration cleanly.
- [x] 1.3 Add `SortOrder int` field to `domain.Verse` struct in `backend/internal/domain/verse.go`. Verify the field is present and the struct compiles.

## 2. BE — Verse usecase: sort_order on create + reorder method

- [x] 2.1 Add `GetMaxSortOrder(ctx context.Context, userID uuid.UUID) (int, error)` to `VerseRepository` interface in `backend/internal/usecase/verse_repository.go` and implement it in `backend/internal/repository/verse_repository.go` with `SELECT COALESCE(MAX(sort_order), -1) FROM verses WHERE user_id = ?`. Verify it compiles.
- [x] 2.2 Update `verseUsecase.Create` in `backend/internal/usecase/verse_usecase.go` to call `GetMaxSortOrder` and set `verse.SortOrder = max + 1` before calling `verseRepo.Create`. Verify: the `fakes_test.go` fake for `VerseRepository` is updated to stub `GetMaxSortOrder`.
- [x] 2.3 Add `ListIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)` to `VerseRepository` interface and implement it in the repository. Verify it compiles.
- [x] 2.4 Add `UpdateSortOrders(ctx context.Context, userID uuid.UUID, orderedIDs []uuid.UUID) error` to `VerseRepository` interface and implement it in the repository. The implementation SHALL loop through `orderedIDs` and issue `UPDATE verses SET sort_order = ? WHERE id = ? AND user_id = ?` for each; wrap in a transaction. Verify it compiles.
- [x] 2.5 Add `Reorder(ctx context.Context, userID uuid.UUID, orderedIDs []uuid.UUID) error` to `VerseUsecase` interface and implement it in `verseUsecase`. The usecase SHALL call `ListIDs` to fetch existing IDs for the user, verify the submitted set is an exact match (returns `ErrVerseReorderIDsMismatch` otherwise), then call `UpdateSortOrders`. Verify: add `ErrVerseReorderIDsMismatch` sentinel to `verse_repository.go` (usecase errors file).
- [x] 2.6 Update `verseRepository.List` in `backend/internal/repository/verse_repository.go` to `ORDER BY sort_order ASC` (replacing the current `last_used_at DESC NULLS LAST, name ASC`). Verify the updated ordering compiles.

## 3. BE — Verse handler: reorder endpoint

- [x] 3.1 Add `reorderVerseRequest` struct (`{ IDs []uuid.UUID json:"ids" validate:"required,min=1" }`) and `Reorder(c echo.Context) error` handler method to `VerseHandler` in `backend/internal/handler/verse_handler.go`. Handler calls `verseUsecase.Reorder` and returns 204 on success, 422 on `ErrVerseReorderIDsMismatch`, 404 on `gorm.ErrRecordNotFound`. Verify it compiles.
- [x] 3.2 Register `POST /verses/reorder` route in `backend/cmd/server/main.go` **before** `PATCH /verses/:id`. Verify routing order in the server file.

## 4. BE — Tests

- [x] 4.1 Add unit tests in `backend/internal/usecase/verse_usecase_test.go`: `Create` assigns `sort_order = max + 1`; `Reorder` succeeds with exact match; `Reorder` returns `ErrVerseReorderIDsMismatch` with partial list; `Reorder` returns `ErrVerseReorderIDsMismatch` with a foreign ID. Update `fakes_test.go` stubs as needed. Verify: `go test ./internal/usecase/...` passes.
- [x] 4.2 Add integration tests in `backend/internal/repository/verse_repository_integration_test.go`: `Create` appends `sort_order` one greater than current max; `List` returns verses ordered by `sort_order ASC`; `UpdateSortOrders` sets positions correctly in a transaction; `GetMaxSortOrder` returns -1 when no verses exist. Verify: integration test suite passes.
- [x] 4.3 Add handler tests in `backend/internal/handler/verse_handler_test.go` for `POST /verses/reorder`: success returns 204; `ErrVerseReorderIDsMismatch` returns 422; unauthorized returns 401. Verify: `go test ./internal/handler/...` passes.

## 5. BE — Bruno file and lint

- [x] 5.1 Create `POST /verses/reorder` Bruno request file following the existing Bruno file convention in the project. Verify the file is syntactically valid Bruno JSON/YAML.
- [x] 5.2 Run `golangci-lint run ./...` from `backend/` and fix any issues introduced by this change. Verify lint exits with code 0.

## 6. FE — Dependencies and verse API hooks

- [x] 6.1 Install `@dnd-kit/core` and `@dnd-kit/sortable` via `npm install` in `frontend/`. Verify both packages appear in `package.json` dependencies.
- [x] 6.2 Add `verse_id: string | null` to the `Character` interface in `frontend/src/features/characters/api/useCharacters.ts`. Verify TypeScript compilation succeeds (no new errors related to this field).
- [x] 6.3 Create `frontend/src/features/verses/api/useVerses.ts`: exports `Verse` interface (`id`, `name`, `sort_order`, `last_used_at: string | null`, `created_at`, `updated_at`), `VERSES_QUERY_KEY`, and `useVerses` query hook hitting `GET /verses`. Verify: hook compiles and `Verse` type is exported.
- [x] 6.4 Create `useCreateVerse.ts`, `useUpdateVerse.ts`, `useDeleteVerse.ts`, and `useReorderVerses.ts` in `frontend/src/features/verses/api/`. Each invalidates `VERSES_QUERY_KEY` on success. `useUpdateVerse` also handles 409 by returning the error so callers can show inline feedback. Verify: all hooks compile.

## 7. FE — VerseCombobox

- [x] 7.1 Create `frontend/src/features/verses/components/VerseCombobox.tsx` mirroring `ArtistCombobox`: same `Combobox` primitives, same clear-X button, props `value: Verse | null`, `onChange`, `placeholder`, `disabled`. Items are sorted by `last_used_at DESC NULLS LAST, name ASC` client-side before rendering. Empty state text: "No verses found". Verify: component compiles with no TypeScript errors.

## 8. FE — VersesModal

- [x] 8.1 Create `frontend/src/features/verses/components/VersesModal.tsx`. The modal renders the verse list fetched from `useVerses`. Each row shows: drag handle (`@dnd-kit/sortable`), verse name, character count (derived from `characters` prop), edit button, delete button. Drag end triggers `useReorderVerses`. Verify: component renders without runtime errors in isolation (can stub dnd context).
- [x] 8.2 Implement inline rename within `VersesModal`: clicking the pencil icon replaces the row with an `<input>` pre-filled with the current name. Enter saves (calls `useUpdateVerse`); Esc cancels. A 409 response from the API is caught and displayed as a small error message beneath the input row (not a toast). Only one row may be in rename mode at a time. Verify: rename flow works end-to-end; duplicate name shows inline error.
- [x] 8.3 Implement inline delete confirmation within `VersesModal`: clicking the trash icon replaces that row with a confirmation prompt showing `"{name}" — {n} characters will move to Unsorted.` and Cancel / Delete buttons. Character count is derived from the `characters` prop. Confirming calls `useDeleteVerse`. Verify: confirmation shows correct character count; deleting removes the verse row.
- [x] 8.4 Implement "+ New verse" button in `VersesModal`: clicking it appends a new inline input row at the bottom of the list (or at the top if the list is empty). Enter saves (calls `useCreateVerse`); Esc cancels. Duplicate name shows inline error. Verify: creating a verse adds it to the list; empty input is rejected.
- [x] 8.5 Implement empty state in `VersesModal` ("No verses yet" heading + "Create one to group characters on the list." subtext) shown when the verse list is empty and no new-verse input is open. Verify: empty state renders correctly when no verses exist.

## 9. FE — CharactersPage grouped layout

- [x] 9.1 Add "Verses" button to the `CharactersPage` header (left of "New character"), and wire it to open `VersesModal`. Add `versesModalOpen` state. Pass `characters` array as prop to `VersesModal` for character count calculation. Verify: button renders and opens the modal.
- [x] 9.2 Add `useVerses` query to `CharactersPage`. Build the grouped layout using `useMemo`: verses in `sort_order` order, each group containing the characters with that `verse_id` that match the current search. Groups with zero matching characters are omitted. "Unsorted" is rendered at the bottom if any characters with `verse_id === null` (or unknown verse) match the search. Verify: characters appear in correct groups; groups with no characters are hidden.
- [x] 9.3 Replace `CharactersGrid` in `CharactersPage` with the grouped view: each section has a `<h2>` or equivalent heading with the verse name and character count, followed by the characters in that group using the existing `ResourceCard` grid layout. The "Unsorted" section uses the same layout with "Unsorted" as heading. Verify: layout matches the mockup (verse name heading + grid of cards + divider between sections).

## 10. FE — CharacterFormModal verse picker

- [x] 10.1 Add `verse` state (`Verse | null`) to `CharacterFormModal`, initialised from `character?.verse_id` (resolved via a `useVerses` lookup). Place `VerseCombobox` above the Notes field. Verify: the combobox renders in the correct position and is pre-filled when editing a character that has a verse.
- [x] 10.2 Update `isDirty()` in `CharacterFormModal` to include verse change detection. Update `resetForm()` to reset verse state. Include `verse_id` in the create and update mutation payloads (`null` to clear, UUID string to assign). Verify: submitting with a verse change calls the API with the correct `verse_id`; the discard-changes dialog fires when only the verse is changed.

## 11. FE — Build and lint

- [x] 11.1 Run `npm run build` from `frontend/` and resolve any TypeScript or build errors introduced by this change. Verify: build exits with code 0.
- [x] 11.2 Run `npm run lint` from `frontend/` and fix any lint issues introduced by this change. Verify: lint exits with code 0.

## Workflow follow-up

- Archive the change once all tasks are complete and the implementation has been reviewed.
