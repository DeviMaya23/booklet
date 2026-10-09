## Context

The app has no meaningful landing page. `/app` redirects to `/app/characters`. Adding a dashboard gives the user an at-a-glance view of active commissions, recent artpieces, and housekeeping gaps, without adding new data models or storage.

Existing list endpoints (`GET /commissions`, `GET /artpieces`) are not reused on the dashboard — they return more data than needed and don't cover the housekeeping queries. A dedicated endpoint keeps the dashboard payload tight and avoids coupling it to the shape of existing list responses.

## Goals / Non-Goals

**Goals:**
- Single `GET /dashboard` endpoint that returns all data the page needs in one round trip
- Dashboard page rendering recent artpieces, in-progress commissions (with mark-contacted), and housekeeping items
- Upload button as a second entry point into the existing Add Files modal
- Dashboard as the default landing page after login

**Non-Goals:**
- Saved filters panel (deferred)
- `+ New Commission` flow (deferred)
- Caching or pagination of dashboard data
- Real-time updates or polling

## Decisions

### Single dedicated endpoint over reusing existing list endpoints

**Decision**: `GET /dashboard` runs its own queries and returns a dashboard-specific shape.

**Rationale**: The in-progress list needs only `waitlist`/`wip` commissions; the housekeeping queries (especially artpieces with no files) require predicates that don't exist on `GET /artpieces`. Filtering client-side from full lists would over-fetch and still require a new endpoint for the artpieces-no-files case. A single endpoint is one network call, one loading state, and a shape the page can consume directly without client-side derivation.

**Alternative considered**: Reuse `useCommissions()` and `useArtpieces()` hooks and filter client-side. Rejected because it loads all statuses and all artpieces on every dashboard visit, and the artpieces-with-no-files check cannot be done client-side (files are not preloaded in the list response).

### Sequential queries over concurrent goroutines

**Decision**: Dashboard usecase runs its ~6 queries sequentially, not with `errgroup`.

**Rationale**: Collection sizes at this scale make the latency difference negligible. Sequential code is easier to read, easier to test, and error propagation is straightforward. Revisit only if dashboard load time becomes a real, observed problem.

### "Artpieces with no files" via NOT EXISTS subquery

**Decision**: The housekeeping query for artpieces with no files uses `NOT EXISTS (SELECT 1 FROM files WHERE artpiece_id = artpieces.id)`.

**Rationale**: The artpiece `List` repo method doesn't preload `Files` (only `CoverFile`). Using `cover_file_id IS NULL` as a proxy would work today (cover is only ever nulled when the last file is removed) but is fragile if a cover-clearing endpoint is added later. The subquery is correct by definition.

### Staleness: null `last_contacted_at` is NOT stale

**Decision**: Follow the same `isStale` logic as the commission list: only contacts older than 14 days go orange. `null` (never contacted) renders in normal muted color with "not contacted yet" text.

**Rationale**: Matches existing user-facing behavior on the Commissions page. Consistency matters more than the theoretical case for treating null as stale.

### Housekeeping scope: all statuses included for "no artist" and "done + no artpieces"

**Decision**: "Commissions with no artist" includes all statuses (waitlist, wip, done). "Done + no artpieces" is scoped to `status = 'done'` only by definition.

**Rationale**: No business rules enforce artist assignment at any status. A done commission with no artist is just as worth flagging as an active one.

### Response shape: flat arrays per category

**Decision**: Dashboard response has top-level keys — `recent_artpieces`, `in_progress`, `housekeeping` (itself an object with four arrays).

```
GET /dashboard → {
  recent_artpieces: [{ id, title, artist_name, thumbnail_url }],   // max 5
  in_progress: [
    { id, title, artist_name, artist_link, status, paid,
      last_contacted_at, created_at }
  ],
  housekeeping: {
    commissions_no_artist:  [{ id, title }],
    artpieces_no_artist:    [{ id, title }],
    done_no_artpieces:      [{ id, title }],
    artpieces_no_files:     [{ id, title }]
  }
}
```

**Rationale**: Clean, predictable, easy to type on the FE. Housekeeping items only need id + title for linking out — no richer data needed.

### Mark-contacted on dashboard: call patch directly, no cell-error abstraction

**Decision**: Dashboard rows call `usePatchCommission` directly rather than going through `useCommissionActions` (which is wired to table-specific `setCellError`).

**Rationale**: `useCommissionActions` is tightly coupled to the commission table's error-per-cell model. Dashboard rows need only `stampLastContacted` — a single `PATCH` with `lastContactedAt`. Extracting only that function is cleaner than adapting the full hook.

### Sidebar position: Dashboard item above Inbox card

**Decision**: A new Dashboard nav item renders above the existing Inbox card in `AppSidebar`. The `/app` redirect changes from `/app/characters` to `/app/dashboard`.

## Risks / Trade-offs

- **Dashboard data staleness after actions** — After mark-contacted on the dashboard, the dashboard query key must be invalidated so the timestamp updates immediately. The `PATCH /commissions/:id` mutation will need to also invalidate `['dashboard']` query key (or the dashboard's own query key). → Mitigation: invalidate `['dashboard']` in `usePatchCommission`'s `onSuccess`, or give the dashboard its own mutation wrapper.

- **Housekeeping subquery performance** — `NOT EXISTS (SELECT 1 FROM files WHERE artpiece_id = artpieces.id)` is unindexed on `files.artpiece_id` if no index exists there. → Mitigation: check if an index exists; add one if not (migration-free, additive-only).

- **`cover_file_id IS NULL` divergence** — If a cover-clearing endpoint is added in the future, the subquery approach remains correct while the proxy would silently break. → No mitigation needed; the subquery is correct by design.
