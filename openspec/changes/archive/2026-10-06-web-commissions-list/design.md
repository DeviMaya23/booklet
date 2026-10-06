## Context

The commissions backend CRUD is fully implemented (create, read, update, delete, artpiece attach/detach). The frontend has no commissions surface at all — no page, no route, no feature folder. This change adds the list screen as the first FE commissions feature, plus the backend additions required to support it (`last_contacted_at` field, `artist_link` in response, and a PATCH endpoint for partial updates).

## Goals / Non-Goals

**Goals:**
- Add `last_contacted_at` to the DB and domain, and expose it in the commission response
- Add `artist_link` to the commission response (field already exists on `Artist`, just not mapped)
- New `PATCH /commissions/:id` endpoint scoped to the four inline-editable fields
- Commissions list page at `/app/commissions` with toggle, sortable columns, and autosave

**Non-Goals:**
- Commission creation modal (deferred)
- Commission detail/edit view (deferred)
- Server-side filtering by status (client-side is sufficient for a single-user list)
- "All commissions" view (deliberately excluded from v1)

## Decisions

### PATCH endpoint fields

The PATCH endpoint accepts only the four fields the list UI edits inline: `status`, `paid`, `paid_date`, `last_contacted_at`. The full PUT remains unchanged for future use by a detail/edit view. Keeping them separate avoids widening the PATCH surface prematurely.

All four fields in the PATCH struct use pointer types (`*string`, `*bool`, `*string` for date, `*time.Time`) so that absent fields are nil and explicit `false`/`null` values are distinguishable from "not provided." The handler only updates fields where the pointer is non-nil.

`status` in PATCH gets the same enum validation as the PUT (`waitlist`, `wip`, `done`). Invalid values return HTTP 422.

### FE orchestrates paid + last_contacted_at together

When the Paid checkbox is checked, the FE sends `{ paid: true, last_contacted_at: now() }` in a single PATCH. The server has no implicit coupling between these fields. This keeps the PATCH semantics simple (what you send is what changes), avoids hidden side effects, and makes the intent explicit in the client.

Unchecking Paid (`paid: false`) does NOT clear `last_contacted_at` — the FE only stamps on transition to `true`.

### Client-side toggle filtering

The toggle between "Waitlist / In Progress" (`status IN ['waitlist', 'wip']`) and "Finished" (`status = 'done'`) is applied client-side after fetching all commissions. A single user's commission list is small enough that this adds no meaningful cost. A server-side filter would require an extra query param and handler change for minimal gain.

### Shadcn components

The following shadcn components are added: `table`, `checkbox`, `calendar`, `popover`. These follow the existing `base-nova` style and are consistent with the project's shadcn setup.

### Relative time display

"10 days ago" for `last_contacted_at` is computed with a small inline utility (no new dependency). The utility formats durations as: "today", "yesterday", "N days ago", "N weeks ago", "N months ago". No `date-fns` needed.

### Autosave error state

Failed saves show an inline error indicator on the affected cell (e.g., red border + icon). Success is silent. This matches the pattern described in the proposal and avoids toast-on-every-save noise.

## Risks / Trade-offs

- **PATCH with pointer types in Go**: The handler must check each field for nil before applying. Missing a nil check silently zeroes a field. Covered by unit tests on each patchable field.
- **FE orchestrating two fields on paid=true**: If the PATCH call fails after the user checks Paid, neither field updates — the failure is atomic at the HTTP level. This is correct behavior; the cell shows an error and the user can retry.
- **No optimistic update**: Cells reflect server state after the mutation resolves. The list re-fetches on mutation success. This is intentional for simplicity; the latency is negligible for a local/same-region API.

## Migration Plan

1. Run migration `000027` (adds `last_contacted_at` column, nullable, no default)
2. Deploy backend (PATCH endpoint, updated response struct)
3. Deploy frontend (new page, components)

Rollback: migration `000027` down drops the column. The frontend can be reverted independently since it's behind a new route with no effect on existing pages.
