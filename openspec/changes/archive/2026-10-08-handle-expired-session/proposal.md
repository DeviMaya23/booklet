## Why

When a user's session expires, API calls silently fail with a 401 and React Query treats it as a generic error — the user gets no feedback and no path back to the app. They need to be redirected to the login page with a clear prompt to log back in.

## What Changes

- `apiFetch` detects 401 responses and signals a new `sessionExpiredStore` (mirroring the `maintenanceStore` pattern)
- A new `sessionExpiredStore` module holds the flag and notifies React subscribers
- `AppShell` reads the store, fires a single toast ("Please log back in"), and redirects to `/`
- 401 responses are deduplicated — only the first one triggers the toast and redirect

## Capabilities

### New Capabilities

- `web-expired-session`: Frontend detection and handling of expired sessions — 401 interception in `apiFetch`, store-driven redirect to the login page, and a single user-facing toast

### Modified Capabilities

- `web-maintenance`: `apiFetch` gains a second side-effect output (session expired signal) alongside the existing maintenance header check — the pattern is unchanged, just extended

## Impact

- `frontend/src/lib/api.ts` — add 401 detection, call `setSessionExpired`
- `frontend/src/lib/sessionExpiredStore.ts` — new file (mirrors `maintenanceStore.ts`)
- `frontend/src/components/AppShell.tsx` — reads store, redirects + toasts once
- No backend changes
- No new dependencies
