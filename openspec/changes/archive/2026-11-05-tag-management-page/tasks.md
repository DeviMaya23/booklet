# Tasks

## 1. Backend: usage count on `GET /tags`

- [x] 1.1 Add `domain.TagWithUsage` (embeds `Tag`, plus `ArtpieceCount int`) in `backend/internal/domain/tag.go`; verify `go build ./...` passes and `domain.Tag` is unchanged
- [x] 1.2 Add `WithUsage bool` (`query:"with_usage"`) to `ListTagFilters`, and `ListWithUsage` to the `TagRepository` interface and `tagUsecase`; verify the build passes
- [x] 1.3 Implement `tagRepository.ListWithUsage`: `LEFT JOIN artpiece_tags`, `WHERE tags.user_id`, `GROUP BY tags.id`, `COUNT(artpiece_tags.artpiece_id)`, qualified ordering shared with `List` (`tags.name`, `tags.last_used_at`); verify existing `List` integration tests still pass with the qualified ordering
- [x] 1.4 Update `ListTags` handler: without `with_usage` return the existing `{id, name}` items unchanged; with it call `ListWithUsage` and return items with an always-present `artpiece_count`; a non-boolean value yields 400 via `c.Bind`. Grep every `GET /tags` / `useTags` / `tagResponse` call site across backend, frontend and extension and confirm none breaks
- [x] 1.5 Usecase unit tests (fakes per CONVENTIONS.md): `ListWithUsage` passes user and filters through and returns the repository result; repository error is propagated and recorded. Handler unit tests: default response has no `artpiece_count`; `with_usage=true` includes it; `with_usage=true&sort=recent` passes both filters; `with_usage=false` equals default; `with_usage=maybe` → 400; assert response bodies, not just status
- [x] 1.6 Repository integration tests for `ListWithUsage`: counts match assignments; unused tag is 0 not absent; an artpiece with several tags counts once per tag; another user's identically named tag and artpieces are not counted; ordering by name; ordering with `sort=recent` (null `last_used_at` last, ties by name); verify `go test` integration suite passes

## 2. Backend: merge endpoint

- [x] 2.1 Add `ErrTagMergeSameTag`, `MergeTagParams{TargetID}`, `ReassignArtpieces(ctx, fromTagID, toTagID)` to the usecase tag interfaces and `Merge(ctx, sourceID, userID, params)` to `tagUsecase`; verify the build passes
- [x] 2.2 Implement `tagRepository.ReassignArtpieces` as `INSERT INTO artpiece_tags (artpiece_id, tag_id) SELECT artpiece_id, ? FROM artpiece_tags WHERE tag_id = ? ON CONFLICT DO NOTHING`, reading the transaction via `dbFromContext`
- [x] 2.3 Implement `tagUsecase.Merge`: reject source == target; inside `Transactor.InTransaction` verify both tags via `GetByIDsAndUserID` (fewer than two → `gorm.ErrRecordNotFound`), call `ReassignArtpieces`, then `Delete(source)`; never touch `last_used_at`. Inject `Transactor` into `NewTagUsecase` and update its call site in `cmd/server/main.go`; verify the build passes
- [x] 2.4 Add `MergeTag` handler (`uuid.Parse` for the path id → 400; `c.Bind` + `c.Validate` on `{target_id}` with `required` → bind error 400 / missing 422; same-tag → 422; `gorm.ErrRecordNotFound` → 404; success → 204) and register `POST /tags/:id/merge` in `cmd/server/main.go`
- [x] 2.5 Usecase unit tests: source == target returns `ErrTagMergeSameTag` and calls nothing else; missing/foreign tag returns `gorm.ErrRecordNotFound` and does not reassign or delete; success calls reassign then delete for the right ids inside the transactor; a reassign failure does not delete the source; a delete failure is returned. Handler unit tests: 204 on success, 400 malformed path id, 400 invalid JSON / non-UUID `target_id`, 422 missing `target_id`, 422 same tag, 404 not found, 500 on unexpected error; assert the exact status and error body
- [x] 2.6 Repository integration tests for `ReassignArtpieces` and the merge flow (the repository the feature depends on): artpieces only on the source gain the target; artpieces with both keep exactly one target row; after the source is deleted its join rows are gone; target `last_used_at` and name unchanged; a foreign user's tag is not matched by `GetByIDsAndUserID`; verify the integration suite passes
- [x] 2.7 Bruno: add `collection/tags/merge_tag.bru` (path `:id`, body `target_id`, docs for 204/400/404/422) and add a `with_usage` query param (disabled by default) with updated docs to `collection/tags/list_tags.bru`; verify both open and run against a local server
- [x] 2.8 Run `golangci-lint run` in `backend/` and fix every issue; verify it exits clean and `git diff go.mod` shows no change

## 3. Frontend: data layer

- [x] 3.1 Add `withUsage` to `UseTagsOptions` in `useTags.ts`, append `'usage'` to the query key only when set, request `?with_usage=true` (combined with `sort`), and add optional `artpiece_count?: number` to `Tag`; extend `useTags.test.tsx` for URL, key separation from non-usage lists, and unchanged behaviour for existing callers. Grep all `useTags(` call sites (`TagPicker`, `ArtistCharacterFilter`, `SaveFilterModal`, others) and verify none changed behaviour
- [x] 3.2 Add `useUpdateTag` (PATCH `/tags/:id`; surface 409 as the existing `TagConflictError`) with a colocated test for request shape, 409 mapping and cache invalidation
- [x] 3.3 Add `useDeleteTag` (DELETE `/tags/:id`) with a colocated test for request shape and invalidation
- [x] 3.4 Add `useMergeTag` (POST `/tags/:id/merge`, body `{target_id}`) with a colocated test for request shape, error mapping and invalidation
- [x] 3.5 In all three mutations invalidate `TAGS_QUERY_KEY` and `ARTPIECES_QUERY_KEY`; verify in `ARTPIECES_QUERY_KEY` usage that per-artpiece detail keys are covered by that prefix (add the detail key if not); verify tests assert both invalidations
- [x] 3.6 Add `features/tags/lib/findTagByName.ts` (trim, case-insensitive, optional excluded id) with a test covering exact, different-case, whitespace, self-match excluded and no match

## 4. Frontend: Tags page

- [x] 4.1 Add the Tags nav item (`Tags` Lucide icon, `/app/tags`) after Artists in `AppSidebar.tsx` and the `/app/tags` route in `router.tsx`; update sidebar tests for MANAGE order (Commissions, Artists, Tags), icon and active state, and verify the link navigates
- [x] 4.2 Build `DeleteTagDialog` (`{ item, onCancel, onConfirm }`, mirroring `DeleteCommissionDialog`) with copy that states the artpiece count (including the zero case) and a colocated test for text, cancel and confirm
- [x] 4.3 Build `MergeTagDialog` (`{ item: {source, target} | null, onCancel, onConfirm }`) with the consequence copy and "This can't be undone" notice, and a colocated test for text, cancel and confirm
- [x] 4.4 Build `TagRow`: name, usage label (singular/plural, link to `/app/artpieces?tags=<id>` only when count > 0), "Unused" badge, edit and delete actions, and inline edit mode (focused pre-filled input, Enter/save, Esc/cancel, hint text, no request when unchanged, blank keeps input open, error on failure); colocated tests for each of those behaviours
- [x] 4.5 Build `TagsList` (rows plus the "No tags yet" empty state, shown only when the user has no tags, not when search matches nothing) with a colocated test
- [x] 4.6 Build `TagsPage`: search box (client-side, case-insensitive substring), sort dropdown ("Recently used" default / "Name (A–Z)") driving `useTags({ sort, withUsage: true })`, one-row-at-a-time `editingId`, rename submit that uses `findTagByName` to route collisions to `MergeTagDialog` (confirm → `useMergeTag`, cancel → stay in edit mode) and everything else to `useUpdateTag`, delete flow via `DeleteTagDialog`; follow the hook-per-concern guidance in CONVENTIONS.md for any state that grows beyond a few lines
- [x] 4.7 Page-level tests: sort and search combine; collision (including different case) opens merge and sends nothing; case-only self rename uses PATCH; merge confirm calls the merge hook with the right ids; delete confirm; empty state; dependent tag data refreshes (a tag picker no longer offers a merged tag)
- [x] 4.8 Run `npm run build` and `npm run lint` in `frontend/` and fix every issue; verify both exit clean and the full frontend test suite passes

## 5. Integration check

- [x] 5.1 Run the app locally and walk the mockup states against real data: list with counts and Unused badges, both sorts, search, rename, rename into an existing name (merge), delete, empty state; verify against the four mockup frames
- [x] 5.2 Run `/check` (OpenSpec verify and FE review) and resolve what it flags

## Workflow follow-up

- Archive the change after the project's review requirements are satisfied.
