# Design

## Context

See `proposal.md` for motivation and scope. Constraints observed in the code:

- **Filtering is client-side.** `filterAndSortArtpieces` filters the full `useArtpieces()` list. The backend's `tag_ids`/`tag_match` list params exist but the frontend does not use them for characters or verses, and this change keeps that pattern.
- **Filter state has three homes** in `useArtpiecesFilter`: the URL (source of truth for render), `localStorage` (`StoredFilter`, restored on mount), and `pendingFilterRef` (so several setters in one event commit one coherent state). Ids from the URL are resolved to objects by looking them up in the query data (`find` + drop nulls).
- **`TokenInput`** is generic with four callers. It only opens the dropdown after typing, its `createFromText` is synchronous (`FolderPicker` uses it), and it commits free text on blur.
- **The saved-filter payload shape is declared in several places** (`FilterPayload` in `parseFilterPayload.ts`, twice inline in `useCreateSavedFilter.ts`, `SaveFilterModal`, `SavedFilterSnapshot` in `ArtpiecesPage`) and consumed by `filterToHref` and `ChangeThumbnailModal`. The backend stores it opaquely, so this is entirely a frontend concern.
- **Create/edit entry points converge**: the artpieces page and both inbox "New artpiece" paths mount the single `ArtpieceFormModal`; edit is `ArtpieceDetailView` only. `useCreateArtpiece` and `useUpdateArtpiece` each have exactly one caller.
- **Backend facts relied on**: `POST /tags` lowercases and trims the name, returns 409 on a case-insensitive duplicate, and imposes no maximum length; `GET /tags?sort=recent` orders by `last_used_at`; `PUT /artpieces/:id` replaces the tag set; `GET /tags` returns only `id` and `name`.

## Goals / Non-Goals

**Goals:**
- Tags can be assigned at create and edit, shown on the detail view, and used as a gallery filter and in saved filters, matching the existing character UX.
- A deleted tag or character id in any stored or linked filter degrades to "ignored", never to an empty result or an error.

**Non-Goals:**
- Tag management UI, per-tag counts, server-side filtering, and tag filtering in the commission and add-to-artpiece modals.
- Pruning dead ids out of the URL, `localStorage` or stored saved filters.
- Refactoring the existing character/verse code paths beyond what the deleted-id fix requires.

## Decisions

### 1. `TagPicker` wraps `TokenInput`, extended with additive optional props

`TagPicker` lives in `features/tags/components/` and renders `TokenInput` with tag data. `TokenInput` gains optional props only; with none of them passed its behavior is byte-for-byte the current behavior, which the existing `TokenInput` tests and the other three callers verify. The additions:

- an opt-in to open the dropdown on focus (listing suggestions before typing). Characters keep today's type-to-open behavior;
- an async create hook (`text -> Promise<item | null>`) plus a predicate saying whether a create row is offered for the current text. The create row is rendered after the suggestions and takes part in arrow-key navigation. A resolved item is added as a chip; `null` leaves the text in place;
- an in-flight guard so a second Enter or click during creation is ignored.

The async hook is a separate prop from `createFromText`, which stays synchronous for `FolderPicker`. Free text is **not** committed on blur for the async path: blur simply clears the text. Unlike `FolderPicker` (where text is the value), an accidental blur must not call `POST /tags`.

`TagPicker` owns the create semantics: it first looks the text up case-insensitively in the loaded tag list and returns the existing tag if found (no request); otherwise `POST /tags`; on 409 it refetches the list and returns the existing tag; on any other error it toasts and returns `null`.

*Alternatives:* (B) a standalone picker duplicates keyboard, blur and dropdown logic that `CONVENTIONS.md` says to avoid duplicating; (C) folding async creation into `createFromText` would change `FolderPicker`'s contract.

### 2. Eager tag creation

A tag is created on `POST /tags` the moment the user confirms it, not when the artpiece form is saved. Cancelling the form therefore leaves the tag behind (it still appears in later suggestions). This keeps the form's save a single request and avoids a multi-step save with partial failure.

*Alternative:* stage names locally and create on save. Rejected: the form would need synthetic ids in chips, a create-then-assign pipeline, and a rollback story.

### 3. One tags query shape everywhere

`features/tags/api/useTags.ts` mirrors `useCharacters` (`TAGS_QUERY_KEY`, optional `sort`). The picker, the gallery filter, the filter-id resolution and `ChangeThumbnailModal` all use `useTags({ sort: 'recent' })`, so they share one cache entry and one request (characters currently fetch twice, sorted and unsorted). `useCreateTag` invalidates `TAGS_QUERY_KEY` (a prefix, so it covers every sort variant). `useCreateArtpiece` and `useUpdateArtpiece` also invalidate it, because the backend bumps `last_used_at` on assignment.

### 4. Gallery filter: parallel state, independent match mode

`StoredFilter`, the `FilterParams` of `filterAndSortArtpieces`, `commitFilter`, the mount restore, and `hasFilterParams` each gain `tagIds` and `tagMatch` alongside the character fields. URL params are `tags` and `tag_match` (`match` stays the character param). Defaults follow characters: absent `tag_match` means `all`; only `any` is written to the URL.

Specifics that are easy to get wrong:
- The mount-restore early-return ("stored filter is empty") must count `tagIds`, or a tag-only stored filter is silently skipped.
- Stored state from before this change has no tag fields: read as `[]` / `all`.
- `pendingFilterRef` keeps the **raw** id list (as it does for characters) so setters that touch other fields do not drop ids that have not resolved yet (for example while tags are loading). The tag setter writes resolved ids only, which prunes dead ids the same way the character setter does.
- `ArtpieceSummary.tags` is a required field (the backend always returns an array); test fixtures that build summaries are updated.

### 5. Deleted ids: resolve and drop, matching what characters already do

The gallery already resolves character ids against `useCharacters()` data and drops the unresolved ones, so a deleted character is a no-op there today. Tags use the same mechanism, plus one guard: while the tags query is pending the tag filter is not applied (the same loading rule the verse filter has), so a linked filter does not flash an empty or wrong result. This is a spec'd behavior in `web-artpieces-gallery` for both tags and characters, with tests; no gallery code change is needed for characters.

If the tags query **fails**, nothing can be resolved, so the URL's tag ids are applied as given (names unknown) instead of silently showing every artpiece under a URL that says "tagged". `ChangeThumbnailModal` does the same for tag and character ids. This is the only case where an unresolved id is kept. Emptying the tag selection also resets the tag match mode to `all`, so a stale `tag_match=any` cannot linger in the URL or storage.

Dead ids are deliberately left in the URL and `localStorage` until the user edits that field. They are harmless once dropped at filter time, and pruning would add state-sync code to a hook that is already delicate.

`ChangeThumbnailModal` is the one place that does not resolve: it passes raw saved-filter character ids (`name: ''`), so a deleted character in `all` mode matches nothing. It is changed to resolve character and tag ids against `useCharacters()` and `useTags()` data (as it already does for artist and verse) and to include the tags query in its loading state. Saved filters are never rewritten when a tag or character is deleted.

### 6. Saved-filter payload: new fields, one declaration

`tag_ids: string[]` and `tag_match_mode: 'all' | 'any'` are added next to the character fields. To avoid a fourth hand-copied declaration, the inline payload types in `useCreateSavedFilter.ts` are replaced with the existing exported `FilterPayload` from `parseFilterPayload.ts`. This reuses an existing type rather than introducing a new abstraction, but it is a change to how those types are organised, so it is flagged here for review.

`parseFilterPayload` defaults absent fields to `[]` / `all`. `filterToHref` writes `tags` when non-empty and `tag_match=any` only for `any`. `SaveFilterModal` saves resolved ids only (it already receives resolved objects), so new saved filters never contain ids that were dead at save time. In `ArtpiecesPage`, `SavedFilterSnapshot` and `snapshotMatchesCurrent` compare the tag id set and tag match mode alongside the character ones, and the tag change handlers clear the snapshot.

### 7. Form integration

`ArtpieceFormModal`: `selectedTags` state, counted in `isDirty()`, cleared in `resetForm()`, sent as `tagIds`. `ArtpieceDetailView`: `editTags` initialised in `handleEnterEdit` from `artpiece.tags`; `unsavedChangesSummary()` compares id sets (the existing character comparison is factored into a small local helper rather than copied) and adds "tags changed"; the save always sends `tag_ids`. In the detail view the Tags field is full width in its own row between the Artist/Characters row and Notes; in the create modal it sits directly below Characters. View mode shows tag chips in the same row as Artist and Characters.

`useCreateArtpiece` and `useUpdateArtpiece` take a required `tagIds: string[]` and always send `tag_ids`, since update replaces the set; their existing tests are updated.

### 8. `ArtistCharacterFilter` gains optional tag props

`tags`, `onTagsChange`, `tagMatch`, `onTagMatchChange`, mirroring the verse and match-mode props: absent means not rendered and no change for the commission and add-to-existing-artpiece callers. The component keeps its name. The tag input there does not offer creation. Summary label, "Clear all" and the "Clear all visible" condition each gain a tag clause.

## Risks / Trade-offs

- [Eager creation leaves orphan tags when a form is cancelled, and there is no tag management UI yet to delete them] → Names are normalised and de-duplicated by the backend so orphans cannot multiply under case variants; tag management is the planned follow-up.
- [A tag name with no length limit server-side can produce a very long chip] → Chips truncate visually via existing chip styles; a limit is a backend decision and out of scope.
- [Raw dead ids remain in the URL and `localStorage`] → Dropped at filter time and pruned when the user edits the tag field; accepted as equivalent to current character behavior.
- [If the tags or characters query fails, `ChangeThumbnailModal` resolves nothing and shows an unfiltered grid] → Same pre-existing behavior as the verse filter today; not worsened for artist or verse, and not addressed here.
- [`ArtpieceSummary.tags` becoming required breaks test fixtures and any code building summaries by hand] → Compile errors surface every site; fixtures are updated in the same change.
- [Adding props to a shared `TokenInput` can regress its four callers] → Props are opt-in, existing `TokenInput` tests stay unchanged and must keep passing, and the picker's new behavior gets its own tests.
- [Saving only resolved ids means a user who saves while tags are still loading drops all tag ids] → The Save Filter modal is opened from the filter popover after the page has loaded; if the tags query has not resolved, the tag chips themselves are absent, so the saved payload matches what the user sees.

## Migration Plan

Frontend-only; ships with the next frontend deploy. No data migration: existing saved filters and `localStorage` filter state have no tag fields and read as no tags. Rollback is a normal revert. Saved filters created with tag fields are ignored by an older frontend, since unknown payload keys are not read.
