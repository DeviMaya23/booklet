## Why

The app currently has no landing page after login — it drops the user into the Characters list. A dashboard gives a single screen where the user can see what needs attention right now: artpieces just added, commissions in flight, and housekeeping items that have quietly gone stale.

## What Changes

- New `GET /dashboard` API endpoint returning a single dashboard-specific payload (recent artpieces, in-progress commissions, housekeeping items)
- New Dashboard page at `/app/dashboard`, replacing the current `/app` → `/app/characters` redirect
- Dashboard nav item added to the sidebar above the Inbox card
- Upload button on the dashboard opens the existing Add Files modal (new entry point, no new upload surface)
- `+ New Commission` button rendered as a no-op placeholder (deferred, consistent with other deferred entry points)
- Saved filters section omitted entirely from this pass (deferred, needs its own design)

## Capabilities

### New Capabilities

- `web-dashboard`: Dashboard page — layout, recent artpieces thumbnail row, In Progress panel with mark-contacted action and stale highlighting, Housekeeping panel with all-clear empty state, Upload button entry point
- `dashboard-data`: `GET /dashboard` endpoint — aggregates recent artpieces, in-progress commissions, and four housekeeping query results into a single response

### Modified Capabilities

- `web-nav-sidebar`: Dashboard nav item added above the Inbox card; `/app` default redirect changed to `/app/dashboard`

## Impact

- **New route**: `GET /dashboard` (protected, user-scoped)
- **New FE route**: `/app/dashboard`
- **New FE page**: `DashboardPage.tsx`
- **New FE feature folder**: `frontend/src/features/dashboard/`
- **Modified**: `AppSidebar.tsx` — new nav item, new ordering
- **Modified**: `App.tsx` — new route, updated redirect
- **Modified**: `backend/cmd/server/main.go` — new route registration
- **New BE files**: `dashboard_handler.go`, `dashboard_usecase.go` (or inline usecase), `dashboard_repository.go`
- **No schema migrations** — all queries are reads over existing tables
- **No changes to existing list endpoints** (`/commissions`, `/artpieces`) — dashboard data comes exclusively from the new endpoint
