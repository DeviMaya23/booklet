## Context

`apiFetch` is a thin, plain async function that handles all API calls. It already has one side-effect: it reads the `X-Booklet-Maintenance` response header and calls `setMaintenanceActive` from `maintenanceStore`. That store uses `useSyncExternalStore` to let React components subscribe reactively without needing context or providers.

`AppLayout` consumes `useMaintenanceActive()` and conditionally renders `MaintenancePage` — the exact pattern we want to replicate for session expiry.

Currently, a 401 response is treated as a generic fetch failure. React Query marks the query/mutation errored but nothing navigates or notifies the user.

## Goals / Non-Goals

**Goals:**
- Redirect to `/` on any 401 from `apiFetch`
- Show a single toast: "Please log back in"
- Deduplicate — only fire once even if multiple in-flight requests fail with 401 simultaneously

**Non-Goals:**
- Silent token refresh / retry (Kinde's `getToken` already attempts this before we call `apiFetch`)
- Distinguishing expiry reasons (expired vs. revoked vs. backend bug)
- Handling 401s on unauthenticated routes (those don't go through `apiFetch`)

## Decisions

### Decision: Mirror maintenanceStore with a new sessionExpiredStore

**Chosen:** Create `frontend/src/lib/sessionExpiredStore.ts` with the same `useSyncExternalStore` pattern.

**Why over alternatives:**
- *QueryClient global onError*: would require every hook to throw a typed error carrying a status code, touching many files. Also doesn't cover mutations without extra config.
- *React context / event emitter*: more setup, no real advantage over the store pattern that's already proven here.
- The store approach keeps `apiFetch` as a plain function (no React dependency) while giving any component reactive access.

### Decision: Consume the store in AppLayout, not AppShell

**Chosen:** Add `useSessionExpired()` in `AppLayout` alongside the existing `useMaintenanceActive()` check.

**Why:** `AppLayout` is already the layer that intercepts authenticated view rendering. Adding the session expiry redirect there is consistent and keeps `AppShell` unchanged.

The handler in `AppLayout`:
1. Reads `useSessionExpired()`
2. When `true`: fires `toast.error('Please log back in')`, calls `navigate('/')`, resets the store flag

**Why reset the flag:** The flag must be cleared before redirect completes; otherwise a future login→return-to-app cycle would immediately redirect again.

### Decision: Deduplication via the store flag itself

The store flag starts `false` and `setSessionExpired` is a no-op if already `true` (same guard as `setMaintenanceActive`). The first 401 sets it; subsequent ones skip. The flag is reset only on redirect. This is sufficient — no additional debounce needed.

## Risks / Trade-offs

- **React Query continues to mark queries as errored** even after redirect. This is acceptable — the user is gone from those pages. On next login they'll land fresh and React Query will refetch.
- **Multiple toasts if somehow called before React re-renders** — mitigated by the store's `if (sessionExpired === active) return` guard, which prevents redundant listener notifications.
- **AppLayout re-render timing** — there's a small window between `setSessionExpired(true)` in `apiFetch` and `AppLayout` reacting. During this window a second 401 may arrive; the store guard handles it.
