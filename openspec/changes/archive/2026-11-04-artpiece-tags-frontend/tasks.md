# Tasks

Test commands run from `frontend/`: `npx vitest run <path>`. Tests follow `TESTING.md`: mock `apiFetch`, render in a fresh `QueryClientProvider` with `retry: false`, query by role/label/text, no class or markup assertions.

## 1. Tags data layer

- [x] 1.1 Add `features/tags/api/useTags.ts` (`Tag` type, `TAGS_QUERY_KEY`, optional `sort`, mirroring `useCharacters`) and `useCreateTag.ts` (`POST /tags`, invalidates `TAGS_QUERY_KEY`, surfaces a 409 distinguishably from other errors). Verify with `useTags.test.tsx` and `useCreateTag.test.tsx`: request path includes `sort=recent`, create posts `{ name }`, create invalidates the tags query, 409 is distinguishable from a generic failure.
- [x] 1.2 Add `tags: { id: string; name: string }[]` to `ArtpieceSummary` (`useArtpieces.ts`) and `ArtpieceDetail` (`useArtpiece.ts`); update every test fixture that builds these types. Verify `npm run build` type-checks and existing tests still pass.
- [x] 1.3 Add required `tagIds: string[]` to `useCreateArtpiece` and `useUpdateArtpiece`, always send `tag_ids`, and invalidate `TAGS_QUERY_KEY` on success. Update `useCreateArtpiece.test.tsx` and `useUpdateArtpiece.test.tsx`: body contains `tag_ids` (including empty array), tags query invalidated after success. Grep confirms the only callers are `ArtpieceFormModal` and `ArtpieceDetailView`.

## 2. TokenInput extension and TagPicker

- [x] 2.1 Extend `components/TokenInput.tsx` with opt-in props: open-on-focus suggestions, async create handler with a create-row predicate (row rendered after suggestions, included in arrow-key navigation), in-flight guard, and no blur-commit on the async path. Verify the existing `TokenInput.test.tsx` passes unchanged, and new cases cover: focus opens suggestions with empty text, create row appears and is keyboard-reachable, async create adds the returned item, `null` result keeps the text, double Enter while pending creates once, blur does not create.
- [x] 2.2 Build `features/tags/components/TagPicker.tsx`: `useTags({ sort: 'recent' })` suggestions in API order excluding selected, case-insensitive existing-name lookup before creating, `POST /tags` via `useCreateTag`, 409 handling (refetch then select existing), error toast on other failures, disabled state. Verify with `TagPicker.test.tsx` covering every scenario in `specs/web-tag-picker/spec.md` (chips and removal, recency order and `sort=recent` request, substring filtering, selected excluded, create row shown/hidden, create success, Enter on exact existing name makes no `POST`, blank text, server error keeps text, 409 selects existing, disabled).

## 3. Create and edit forms

- [x] 3.1 Add the Tags field to `ArtpieceFormModal` directly below Characters: `selectedTags` state, included in `isDirty()`, cleared in `resetForm()`, passed as `tagIds`. Verify with an `ArtpieceFormModal.tags.test.tsx` (next to the existing recency test): tag selection sends `tag_ids`, no tags sends `[]`, inline-created tag is submitted, tags-only form triggers the discard confirmation on Escape, tags cleared on reopen, save failure keeps selected tags, tags request uses `sort=recent`, tags query invalidated after save.
- [x] 3.2 Add tags to `ArtpieceDetailView`: chips in view mode (or "—"), full-width Tags field in edit mode between the Artist/Characters row and Notes, `editTags` initialised in `handleEnterEdit`, id-set comparison helper shared with the character comparison, "tags changed" in `unsavedChangesSummary()`, `tag_ids` always sent on save, tags query invalidated, Cancel discards. Verify by extending `ArtpieceDetailView.test.tsx` and `ArtpieceDetailView.recency.test.tsx` for every scenario in `specs/web-artpiece-detail/spec.md` (view chips and empty state, init from artpiece, add/create/remove, save with changed/cleared/unchanged tags, caption "tags changed", remove-then-re-add shows no change, blocker on tag-only edit, cancel reverts, `sort=recent`).

## 4. Saved-filter payload and shared filter logic

- [x] 4.1 Extend `FilterPayload` in `lib/parseFilterPayload.ts` with `tag_ids` and `tag_match_mode` (absent reads as `[]` / `all`, non-string ids dropped) and replace the inline payload types in `useCreateSavedFilter.ts` with it. Verify with `parseFilterPayload.test.ts`: tags parsed, absent fields default, malformed `tag_ids` entries dropped; `npm run build` type-checks the replaced types.
- [x] 4.2 Update `lib/filterToHref.ts` to write `tags` when non-empty and `tag_match=any` only for `any`. Verify with `filterToHref.test.ts`: any mode, all mode, no tag fields, invalid payload still returns `/app/artpieces`.
- [x] 4.3 Add `selectedTags` / `tagMatch` (optional, defaulting to none / `all`) to `filterAndSortArtpieces` with all/any semantics over `ArtpieceSummary.tags`. Verify with `useArtpiecesFilter.test.ts`: all, any, single tag, no tags, combined with artist and with characters using different match modes.

## 5. Gallery tag filter

- [x] 5.1 Extend `useArtpiecesFilter`: `tagIds` / `tagMatch` in `StoredFilter` and `pendingFilterRef`, `tags` / `tag_match` URL read/write in `commitFilter`, `hasFilterParams`, restore-from-storage (tag-only stored filter is restorable, pre-tags stored state reads as none), resolve ids against `useTags({ sort: 'recent' })` and drop unresolved, skip the tag filter while the tags query is pending, expose `selectedTags`, `setSelectedTags`, `tagMatch`, `setTagMatch`. Verify with `useArtpiecesFilter.test.ts` for every scenario in the tag-persistence and filter-id-resolution requirements of `specs/web-artpieces-gallery/spec.md`, including deleted tag in `all` mode, every tag deleted, deleted character in `all` mode, tags still loading, and "Clear all" removing both params.
- [x] 5.2 Add the optional tag props (`tags`, `onTagsChange`, `tagMatch`, `onTagMatchChange`) to `components/ArtistCharacterFilter.tsx`: Tags token input (no creation, `sort=recent`), independent All/Any toggle at 2+ tags, summary label clause, Clear all clause and visibility condition. Verify with `ArtistCharacterFilter.test.tsx` and `ArtistCharacterFilter.recency.test.tsx` for every scenario in `specs/web-artist-character-filter/spec.md`, and confirm the commission-modal and add-to-existing-artpiece tests render no Tags field.
- [x] 5.3 Wire `pages/ArtpiecesPage.tsx`: pass tag state to `ArtistCharacterFilter`, include tags in `hasActiveFilter`, extend `SavedFilterSnapshot` / `snapshotMatchesCurrent`, and add tag change handlers that clear the snapshot. Verify with an `ArtpiecesPage` test that selecting a tag filters the grid, the "Save to dashboard" button appears for a tag-only filter, and the "Saved as" confirmation clears when tags or the tag match mode change.

## 6. Save filter modal and dashboard

- [x] 6.1 Update `SaveFilterModal.tsx`: `filterState` gains `tags` and `tagMatch`, the summary chip lists tags and their mode, and the payload sends `tag_ids` (resolved ids only) and `tag_match_mode`. Verify by extending `SaveFilterModal.test.tsx` for the chip, payload with tags, payload with none (`tag_ids: []`), and resolved-only ids.
- [x] 6.2 Fix and extend `ChangeThumbnailModal.tsx`: resolve character and tag ids against `useCharacters()` / `useTags()` data (dropping unresolved), apply the tag filter, include the tags query in the loading state, and add tags to the filter chip label. Verify by extending `ChangeThumbnailModal.test.tsx`: tag filter all and any, deleted character in `all` mode does not empty the grid, deleted tag in `all` mode does not empty the grid, loading state while tags or characters load.
- [x] 6.3 Verify dashboard navigation end to end: extend the dashboard test (or add one beside `filterToHref.test.ts` usage) so a saved-filter tile with tags navigates with `tags` / `tag_match`, a tile without tag fields adds no tag params, and a tile referencing a deleted tag opens the gallery with that tag ignored.

## 7. Integration checks

- [x] 7.1 Run `npx vitest run` from `frontend/` and verify the whole suite passes, including the unchanged `TokenInput`, commission-modal and add-to-existing-artpiece tests.
- [x] 7.2 Run `npm run build` and `npm run lint` from `frontend/` and fix whatever they report; verify both exit 0 (no new `eslint-disable` comments).
- [x] 7.3 Manual check against a running app: create an artpiece from the artpieces page and from the inbox with an existing tag and an inline-created tag; edit tags on the detail page; filter the gallery by tags in All and Any; save the filter, delete one of its tags, and confirm the dashboard tile and the Change thumbnail modal still show results.

## Workflow follow-up

- Archive the change after review requirements are satisfied.
