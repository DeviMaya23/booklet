## Why

Booklet has a complete backend but no frontend yet. This change scaffolds the minimum viable web app: a public home page that drives users into Kinde login, a protected shell page they land on after authentication, and the server-driven maintenance mode that lets the backend gate all traffic during deploys.

## What Changes

- Add `App.tsx` with the full route tree (`/`, `/callback`, `/app` behind `AuthGuard`)
- Add `HomePage` — public landing page with a login button that calls Kinde's login flow
- Add `AppPage` — empty protected shell, the entry point for all future app features
- Clean up all `bookleaf` references in the copied frontend files (storage keys, header names, UI strings, log prefixes, HTML title)
- Wire `useMaintenanceActive()` into the app layout so the maintenance page short-circuits rendering when the backend signals downtime
- Narrow theme system to a single default theme (`warm`); the type and storage key are booklet-specific from here

## Capabilities

### New Capabilities

- `web-home`: Public landing page — unauthenticated entry point with a login trigger
- `web-app-shell`: Protected app shell — the authenticated entry point, empty placeholder for future features
- `web-maintenance`: Server-driven maintenance mode — backend signals downtime via a response header; all authenticated views collapse to a maintenance page; bypass available via a localStorage token

### Modified Capabilities

## Impact

- New files: `frontend/src/App.tsx`, `frontend/src/pages/HomePage.tsx`, `frontend/src/pages/AppPage.tsx`
- Modified files: `frontend/src/main.tsx`, `frontend/src/hooks/useTheme.tsx`, `frontend/src/hooks/useTheme.test.tsx`, `frontend/src/lib/api.ts`, `frontend/src/components/MaintenancePage.tsx`, `frontend/src/components/SimplePageLayout.tsx`, `frontend/index.html`, `frontend/.env.example`
- No backend changes
- No new dependencies — all packages already present in `package.json`
