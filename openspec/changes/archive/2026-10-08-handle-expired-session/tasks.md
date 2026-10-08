## 1. Session Expired Store

- [x] 1.1 Create `frontend/src/lib/sessionExpiredStore.ts` mirroring `maintenanceStore.ts` — module-level flag, `subscribe`/`getSnapshot`, `setSessionExpired`, `useSessionExpired`

## 2. apiFetch — 401 Detection

- [x] 2.1 In `frontend/src/lib/api.ts`, after the `setMaintenanceActive` call, add: `if (res.status === 401) setSessionExpired(true)`

## 3. AppLayout — Redirect on Expiry

- [x] 3.1 In `frontend/src/components/AppLayout.tsx`, read `useSessionExpired()` and `useNavigate()`
- [x] 3.2 Add a `useEffect` that fires when `sessionExpired` is `true`: call `setSessionExpired(false)`, fire `toast.error('Please log back in')`, navigate to `/`

## 4. Tests

- [x] 4.1 Unit test `sessionExpiredStore`: verify `setSessionExpired(true)` notifies listeners, second call is a no-op, and `setSessionExpired(false)` resets
- [x] 4.2 Unit test `apiFetch`: verify a 401 response calls `setSessionExpired(true)` (mock the store)
- [x] 4.3 Unit test `AppLayout`: verify that when `sessionExpired` is `true`, toast fires once, navigate is called with `'/'`, and store is reset

## 5. Build & Lint

- [x] 5.1 Run `npm run build` from `frontend/` and fix any type errors
- [x] 5.2 Run `npm run lint` from `frontend/` and fix any lint issues
