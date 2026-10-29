# Design

## Context

The verse backend (CRUD, `verse_id` on characters, `last_used_at` tracking) is fully implemented across `000033`/`000034` migrations, `verse_handler`, `verse_usecase`, and `verse_repository`. The current `LIST /verses` orders by `last_used_at DESC NULLS LAST, name ASC` — an implicit recency sort the user cannot control.

The character list page is a flat `CharactersGrid`. The `Character` API response already includes `verse_id`.

## Goals / Non-Goals

**Goals:**
- Add `sort_order` to `verses` and a `POST /verses/reorder` endpoint so the user controls group order in the UI.
- Character list page renders characters grouped by verse in `sort_order` order; unsorted characters at the bottom; both views hide empty groups.
- `VersesModal` lets the user create, rename (inline), delete (inline confirmation + count), and drag-to-reorder verses.
- `VerseCombobox` in `CharacterFormModal` mirrors `ArtistCombobox` for assigning/clearing a verse.

**Non-Goals:**
- Filtering characters by verse via the backend `?verse_id=` param (grouping is client-side).
- Public verse visibility or sharing.
- Ordering within a verse group (characters remain ordered by `created_at DESC` as before).

## Decisions

### D1 — Integer sort_order with full-array reorder endpoint

**Chosen:** Simple `INT sort_order` column; new `POST /verses/reorder` accepts `{ "ids": [...] }` (the complete ordered list) and writes `sort_order = index` for every verse in a single transaction.

**Why not fractional indexing:** Solves concurrent multi-user edits and scale — neither applies here. Adds library dependency or careful precision arithmetic for no benefit.

**Why not per-row PATCH:** The FE sends the final state after drag ends, not a delta. A full-array write is atomic, unambiguous, and avoids stale intermediate states.

**Validation:** The usecase fetches all verse IDs for the user and rejects the request (422) if the submitted set is not an exact match — no missing, no extras. This prevents partial reorders that leave some verses with stale `sort_order`.

**On create:** usecase queries `SELECT COALESCE(MAX(sort_order), -1) + 1 WHERE user_id = ?` and sets `sort_order` for the new verse. Single-user app: no meaningful race condition.

**Migration backfill:** `000035` adds the column with `DEFAULT 0`, then backfills distinct values using `ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at)` so existing verses get deterministic stable positions.

**Echo router note:** `POST /verses/reorder` is registered before `PATCH /verses/:id`. Echo's radix tree gives static segments precedence over parameters, so this is safe regardless of registration order — but register static first as convention.

### D2 — Client-side grouping; two independent queries

**Chosen:** `CharactersPage` issues two queries: `useVerses` (sorted by `sort_order`) and `useCharacters` (all characters, unchanged). Grouping is a `useMemo` over the combined data.

```
groups = verses
  .map(v => ({ verse: v, characters: characters.filter(c => c.verse_id === v.id && matchesSearch(c)) }))
  .filter(g => g.characters.length > 0)

unsorted = characters.filter(c => (c.verse_id === null || !verseById[c.verse_id]) && matchesSearch(c))
```

Groups with zero matching characters are not rendered. "Unsorted" is skipped if empty. Search keeps groups visible — only hides groups that have no matching members.

**Why not server-side grouping / `?verse_id=` filter per verse:** Would require `N+1` requests or a new aggregation endpoint. Client-side is trivial given the data is already fetched.

**Defensive handling of stale verse_id:** A character whose `verse_id` doesn't match any verse in the list (shouldn't happen with `ON DELETE SET NULL`, but defensively handled) falls into Unsorted.

### D3 — `GET /verses` ordering changed; picker sorts client-side by `last_used_at`

`LIST /verses` changes from `last_used_at DESC NULLS LAST, name ASC` to `sort_order ASC`. The `VerseCombobox` (verse picker in `CharacterFormModal`) wants most-recently-used first. Since all verse data is already in the React Query cache, the combobox component sorts its local copy by `last_used_at DESC NULLS LAST, name ASC` before rendering items — no extra endpoint or query param.

### D4 — `@dnd-kit/sortable` for drag-to-reorder in `VersesModal`

No drag precedent in the codebase. `@dnd-kit/core` + `@dnd-kit/sortable` chosen over:
- Native HTML5 drag: poor mobile UX, quirky cross-browser scroll behaviour.
- `react-beautiful-dnd`: deprecated, no active maintenance.
- `react-sortable-hoc`: also abandoned.

`@dnd-kit` is lightweight, actively maintained, accessibility-aware, and has no React peer-dep conflicts with the current setup.

Drag ends → optimistic update to local verse order → call `POST /verses/reorder` → on error, React Query cache refetch reverts.

### D5 — Inline rename and inline delete confirmation (no toast, no AlertDialog)

**Rename:** clicking the pencil on a verse row replaces the name with an `<input>`, pre-filled with the current name. Enter saves; Esc cancels. API 409 (name conflict) is shown as a small red message directly under the input row — not a toast — so the user knows what to fix without losing the input state.

**Delete:** clicking the trash icon on a verse row replaces the row with an inline confirmation prompt: _"Delete '{name}'? {n} characters will move to Unsorted."_ + Cancel / Delete buttons. Character count is derived from the `useCharacters` cache. No separate `AlertDialog` popup.

This keeps all verse management interactions within the single `VersesModal` surface, matching the mockup exactly.

### D6 — `VerseCombobox` mirrors `ArtistCombobox`

`VerseCombobox` is a new component in `features/verses/components/` using the same `Combobox` primitives and same clear-X pattern. The only differences from `ArtistCombobox`: data source (`useVerses`, sorted client-side by `last_used_at`), the typed value is `Verse | null`, and empty-state text reads "No verses found".

Placed above the Notes field in `CharacterFormModal`, with `verse` state added to `isDirty()` and `resetForm()`.

## Risks / Trade-offs

- **`sort_order = 0` for all existing verses after migration** → mitigated by the ROW_NUMBER backfill in `000035`. After the migration runs, existing verses have distinct, stable positions.
- **Optimistic reorder + API failure** → the drag's optimistic update is reverted by React Query's `onError` → `queryClient.invalidateQueries(VERSES_QUERY_KEY)`. Visible as a snap-back.
- **VersesModal open while another tab reorders** → single-user app; not a real concern. The `useVerses` query has a standard stale time, so re-opening the modal or returning from background will refetch.
- **`LIST /verses` order change is a breaking API change** for any caller relying on `last_used_at` order. Only the FE combobox used that order; it is updated to sort client-side.

## Migration Plan

1. Deploy `000035_add_sort_order_to_verses` (additive; `DEFAULT 0` is safe on live table, then backfill in same migration).
2. Deploy updated backend with changed `List` ordering and new `/verses/reorder` route.
3. Deploy updated frontend. No feature flag needed — character list grouping is purely additive; the page degrades gracefully if verses is empty.
4. Rollback: `000035` down migration drops the column; backend reverts `List` order. Frontend would show flat list again.
