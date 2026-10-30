# Tasks

## 1. DB Migration + Domain Model

- [x] 1.1 Create a migration that adds `is_reference boolean NOT NULL DEFAULT false` to the `artpieces` table; verify it applies cleanly and all existing rows have is_reference = false
- [x] 1.2 Add `IsReference bool` field to `domain.Artpiece` with GORM tag `gorm:"column:is_reference;not null;default:false"`; verify the project builds

## 2. Backend — CRUD Params, Handler, Response

- [x] 2.1 Add `IsReference bool` to `CreateArtpieceParams` and `*bool` to `UpdateArtpieceParams`; add `IsReference *bool json:"is_reference"` to `createArtpieceRequest` and `updateArtpieceRequest` in the handler; verify the project builds
- [x] 2.2 Pass `IsReference` from `createArtpieceRequest` to `CreateArtpieceParams` in `CreateArtpiece` handler; set it on the `domain.Artpiece` before calling `artpieceRepo.Create`; verify is_reference is persisted in the create integration test
- [x] 2.3 Pass `IsReference` from `updateArtpieceRequest` to `UpdateArtpieceParams` in `UpdateArtpiece` handler; in `artpieceRepository.Update`, add `if params.IsReference != nil { updates["is_reference"] = *params.IsReference }` to the updates map; verify the field is updated correctly and omitting it leaves it unchanged in the integration test
- [x] 2.4 Add `IsReference bool json:"is_reference"` (no omitempty) to `artpieceResponse`; populate it from `a.IsReference` in `toArtpieceResponse`; verify the field appears as `false` and `true` in handler unit test responses
- [x] 2.5 Add/update handler unit tests in `artpiece_handler_test.go`: create with is_reference true returns 201 with is_reference: true; update with is_reference false changes the field; update omitting is_reference does not change the field; response always includes is_reference
- [x] 2.6 Add/update integration tests in `artpiece_repository_integration_test.go`: create artpiece with is_reference true persists correctly; update is_reference from true to false; update omitting is_reference leaves it unchanged; list without filter returns both kinds

## 3. Backend — List Filter

- [x] 3.1 Add `IsReference *bool \`query:"is_reference"\`` to `ListArtpieceFilters`; add `if filters.IsReference != nil { q = q.Where("artpieces.is_reference = ?", *filters.IsReference) }` to `artpieceRepository.List`; verify the project builds
- [x] 3.2 Add integration tests for the list filter: is_reference=true returns only references; is_reference=false returns only non-references; absent param returns all

## 4. Bruno Files

- [x] 4.1 Update `collection/artpieces/create_artpiece.bru` to include `is_reference` in the request body; update `collection/artpieces/update_artpiece.bru` similarly; update `collection/artpieces/list_artpieces.bru` to show the optional `is_reference` query param; verify the files are valid Bruno syntax

## 5. Backend Lint

- [x] 5.1 Run `golangci-lint run ./...` from `backend/`; fix any issues introduced by this change

## 6. Frontend — Types and API Hooks

- [x] 6.1 Add `is_reference: boolean` to `ArtpieceSummary` in `useArtpieces.ts`; add `isReference?: boolean` to `CreateArtpieceParams` in `useCreateArtpiece.ts` and include it in the POST body; add `isReference?: boolean` to `UpdateArtpieceParams` in `useUpdateArtpiece.ts` and include it in the PUT body; verify the project builds (`npm run build`)
- [x] 6.2 Add `is_reference: boolean` to `FilterPayload` in `parseFilterPayload.ts`; coerce absent/non-boolean values with `const is_reference = p.is_reference === true`; add to `CreateSavedFilterParams.filterPayload` in `useCreateSavedFilter.ts`; verify the project builds
- [x] 6.3 Add tests to `parseFilterPayload.test.ts`: is_reference true round-trips; is_reference false round-trips; absent is_reference coerces to false; non-boolean is_reference coerces to false; verify all tests pass

## 7. Frontend — useArtpiecesFilter

- [x] 7.1 Add `isReference` state to `useArtpiecesFilter`: read from URL param `ref` (absent → false, `'true'` → true); expose `setIsReference` that writes `ref=true` or deletes the param; add the filter to `filterAndSortArtpieces` (when false, exclude artpieces with is_reference === true; when true, exclude artpieces with is_reference === false); verify the project builds
- [x] 7.2 Add tests to `useArtpiecesFilter.test.ts`: toggle off hides references; toggle on shows only references; absent ref param defaults to false; URL param `ref=true` restores toggle; other filters still apply when toggle is on; verify all tests pass

## 8. Frontend — ArtistCharacterFilter

- [x] 8.1 Add optional props `isReference?: boolean` and `onIsReferenceChange?: (v: boolean) => void` to `ArtistCharacterFilter`; render a "References only" toggle (Switch or Checkbox) in the panel when both props are present; verify the component renders without them (no regression)
- [x] 8.2 Update `buildSummaryLabel` to include `"references"` when `isReference` is true (e.g. append `· references` to existing parts or return `"references"` when no other parts); verify label in unit/snapshot test
- [x] 8.3 Update the "Clear all" click handler to also call `onIsReferenceChange(false)` when the prop is present; verify the existing "Clear all" behaviour is unchanged when the prop is absent; verify all existing ArtistCharacterFilter tests still pass

## 9. Frontend — ArtpiecesPage (Gallery)

- [x] 9.1 Destructure `isReference, setIsReference` from `useArtpiecesFilter`; pass as `isReference` / `onIsReferenceChange` props to `ArtistCharacterFilter`; verify the toggle appears in the filter popover and filtering works end-to-end
- [x] 9.2 Add `isReference: boolean` to `SavedFilterSnapshot`; include it in `handleSaveFilterSuccess`; add `handleIsReferenceChange` that calls `setSavedFilterSnapshot(null)` then `setIsReference`; pass it as `onIsReferenceChange` instead of the bare setter; verify the "Saved as" confirmation clears when the toggle changes
- [x] 9.3 Add `isReference` to `SaveFilterModal`'s `filterState` prop call; add `is_reference: isReference` to `filterPayload` passed to `useCreateSavedFilter`; verify a saved filter created with the toggle on stores is_reference: true

## 10. Frontend — filterToHref (DashboardPage)

- [x] 10.1 In `filterToHref`, add `if (payload.is_reference) params.set('ref', 'true')` before building the query string; verify that a saved filter with is_reference: true links to `/app/artpieces?ref=true` and one with false does not add the param

## 11. Frontend — ArtpieceFormModal (Create)

- [x] 11.1 Add `isReference` local state (default false) to `ArtpieceFormModal`; render a "This is a reference" checkbox below the notes field; pass `isReference` to `useCreateArtpiece` call; reset the checkbox in `resetForm`; verify the checkbox defaults unchecked, submits the flag, and resets on close

## 12. Frontend — ArtpieceDetailView (Edit)

- [x] 12.1 Add `editIsReference` state to `ArtpieceDetailView`; initialize it from `artpiece.is_reference` in `handleEnterEdit`; render a "This is a reference" checkbox in edit mode; pass `isReference: editIsReference` to `useUpdateArtpiece` call; include the checkbox in `unsavedChangesSummary` dirty check; verify the checkbox reflects the artpiece state, saves correctly, and triggers the unsaved-changes blocker

## 13. Frontend Build and Lint

- [x] 13.1 Run `npm run build` from `frontend/`; fix any type errors or build failures introduced by this change
- [x] 13.2 Run `npm run lint` from `frontend/`; fix any lint issues introduced by this change

## Workflow follow-up

- Archive the change once the implementation has been reviewed and merged.
