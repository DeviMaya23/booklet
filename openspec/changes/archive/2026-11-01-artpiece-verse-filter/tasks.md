# Tasks

## 1. Filter hook — verse dimension

- [x] 1.1 Add `verseId: string | null` to `StoredFilter` and `FilterParams` in `frontend/src/features/artpieces/hooks/useArtpiecesFilter.ts`; build the `characterId -> verse_id` lookup from `useCharacters()` data and extend `filterAndSortArtpieces` to match an artpiece when any of its `characters` resolves to the selected verse. Verify with unit tests in `useArtpiecesFilter.test.ts` (or equivalent): matches an artpiece with a character in the selected verse, excludes one whose characters are all verse-less or in other verses, matches when combined with an active artist/character filter (AND), and resolves to inactive when the stored `verseId` doesn't exist in current `useVerses()` data.
- [x] 1.2 Wire `verseId` into the `?verse=<uuid>` URL param and `commitFilter`/`setSelectedVerse` following the existing `artistId` pattern; verify selecting/clearing a verse updates the URL and `localStorage` under `artpieces.filterState` (existing test file for this hook, extended).

## 2. Shared `ArtistCharacterFilter` component

- [x] 2.1 Add optional `selectedVerse`/`onVerseChange` props to `frontend/src/components/ArtistCharacterFilter.tsx`, rendering a `VerseCombobox` only when both are provided; include verse reset in the "Clear all" handler when `onVerseChange` is provided. Verify with component tests: verse field renders when props given, is absent when omitted (confirming `CommissionFormModal`/`AddToExistingArtpieceModal` usages are unaffected since they won't pass the new props), selecting/clearing calls `onVerseChange` correctly, and "Clear all" resets verse when the prop is present.

## 3. Gallery page wiring

- [x] 3.1 In `frontend/src/pages/ArtpiecesPage.tsx`, pass `selectedVerse`/`setSelectedVerse` from the filter hook into the gallery's `ArtistCharacterFilter` usage, and include verse in the filter-button summary label and the "Save to dashboard" active-filter check. Verify manually in the dev server: selecting a verse filters the grid, the Filter button summary reflects it (e.g. "1 verse"), and it combines correctly (AND) with artist/character/references-only filters already active.

## 4. Saved filter payload — parse, create, update

- [x] 4.1 Add `verse_id: string | null` to `FilterPayload` in `frontend/src/features/artpieces/lib/parseFilterPayload.ts`, defaulting to `null` when absent. Verify with unit tests: parses a present `verse_id`, defaults to `null` when the key is missing (old saved filters), and rejects/ignores a non-string non-null value the same way other fields do.
- [x] 4.2 Add `verse_id` to the `filterPayload` shape in `CreateSavedFilterParams`/`SavedFilter` in `frontend/src/features/artpieces/api/useCreateSavedFilter.ts`.
- [x] 4.3 In `frontend/src/features/artpieces/components/SaveFilterModal.tsx`, include `verse_id` in the payload built by `handleSave`, add the `verse` field to the `filterState` prop type, and extend `buildFilterSummary` to include the verse name when active. Verify with component tests: saving with an active verse filter sends `verse_id` in the create request body, and the summary chip includes the verse name.

## 5. Change-thumbnail picker

- [x] 5.1 In `frontend/src/features/artpieces/components/ChangeThumbnailModal.tsx`, add the `useCharacters()` query (not currently imported), pass the parsed `verse_id` through to `filterAndSortArtpieces`, and extend `filterChipLabel` to mention the verse when active. Verify with component tests: the candidate grid excludes artpieces with no character in the saved filter's verse, and the chip label includes "verse" text when `payload.verse_id` is set.

## 6. Dashboard saved-filter navigation

- [x] 6.1 In `frontend/src/pages/DashboardPage.tsx`, extend `filterToHref()` to map `payload.verse_id` to a `verse` URL param. Verify with unit tests: a payload with a non-null `verse_id` produces an href containing `verse=<id>`, and a payload with `verse_id` null/absent produces no `verse` param.

## 7. Final verification

- [x] 7.1 Run `npm run build` in `frontend/` and fix any type errors introduced by the new `verse_id`/`selectedVerse` fields across the touched files.
- [x] 7.2 Run `npm run lint` in `frontend/` and fix any findings.
- [x] 7.3 Run `npm run test` in `frontend/` and confirm all new and existing tests pass.

## Workflow follow-up

- Archive the change after review.
