# Tasks

## 1. Core implementation in `useArtpiecesFilter`

- [x] 1.1 Add `StoredFilter` interface (`{ artistId, characterIds, matchMode, isReference }`) and `readStoredFilter` / `writeStoredFilter` helpers (localStorage key `artpieces.filterState`, reads/writes wrapped in try/catch). Verify the helpers compile and the key name matches the design.
- [x] 1.2 Add a mount-only `useEffect` (empty dep array) that, when the URL has no filter params, reads localStorage and calls `setSearchParams` with the stored values if they are non-trivial (artist or characters or isReference set). Verify by adding the scenario test in task 2.2 and confirming it passes.
- [x] 1.3 In the same mount effect, handle the URL-params-present path: when filter params are already in the URL (e.g. arrived from a saved-filter link), call `writeStoredFilter` with the current URL state. Verify by adding the scenario test in task 2.3 and confirming it passes.
- [x] 1.4 Call `writeStoredFilter` in each structural filter setter — `setSelectedArtist`, `setSelectedCharacters`, `setCharacterMatch`, `setIsReference` — writing the full current filter state with the new value applied. Verify by adding the setter tests in task 2.4 and confirming they pass.

- [x] 1.5 Route every setter through a single `commitFilter` that writes localStorage and serializes all four filter params in one non-updater `setSearchParams` call, built from `pendingFilterRef` on top of a `paramsRef` snapshot so non-filter params are preserved. Verify with the multi-setter tests in task 2.7.

## 2. Unit tests

- [x] 2.1 Add `localStorage` mock setup to the `beforeEach` in `useArtpiecesFilter.test.ts` (reset `localStorage` between tests using `vi.stubGlobal` or direct mock). Verify all existing tests still pass after the mock is added.
- [x] 2.2 Add test: mount with no URL params and a non-trivial value in localStorage → `setSearchParams` is called to restore the stored artist/characters/ref. Verify the test passes.
- [x] 2.3 Add test: mount with URL filter params present → `localStorage.setItem` is called with the URL-derived state. Verify the test passes.
- [x] 2.4 Add tests for each setter: `setSelectedArtist`, `setSelectedCharacters`, `setCharacterMatch`, `setIsReference` — each call writes the updated full filter state to `localStorage`. Verify the tests pass.
- [x] 2.5 Add test: mount with no URL params and no localStorage data → `setSearchParams` is NOT called (no spurious update). Verify the test passes.
- [x] 2.6 Add test: `setSelectedArtist(null)` with no characters or isReference active → localStorage is written with an empty state (all fields null/empty/false), so the next visit starts blank. Verify the test passes.
- [x] 2.7 Make the `setSearchParams` mock faithful to React Router (the updater form receives the params of the last render, so it does not chain), then add tests: "Clear all" (artist + characters + ref setters in one `act`) clears every filter param and writes empty state to localStorage; two setters in one event both take effect; a filter change preserves unrelated search params. Verify the tests fail against a per-setter delta implementation and pass after task 1.5.

## 3. Build and lint

- [x] 3.1 Run `npm run lint` from `frontend/` and fix any issues.
- [x] 3.2 Run `npm run build` from `frontend/` and fix any issues.

## Workflow follow-up

- Archive the change once the PR is reviewed and merged.
