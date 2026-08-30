## 1. bookleaf → booklet rename sweep

- [x] 1.1 Update `src/main.tsx`: replace `[bookleaf]` log prefix with `[booklet]`
- [x] 1.2 Update `index.html`: change `<title>` to "Booklet" and update the inline theme-restore script's localStorage key to `booklet-theme`
- [x] 1.3 Update `.env.example`: replace "Bookleaf" in the API base URL comment
- [x] 1.4 Update `src/hooks/useTheme.tsx`: rename storage key to `booklet-theme` and narrow `Theme` type to `'warm'` only
- [x] 1.5 Update `src/hooks/useTheme.test.tsx`: replace all `bookleaf-theme` localStorage key references with `booklet-theme`
- [x] 1.6 Update `src/lib/api.ts`: rename `bookleaf-maintenance-bypass` → `booklet-maintenance-bypass`, `X-Bookleaf-Bypass` → `X-Booklet-Bypass`, `X-Bookleaf-Maintenance` → `X-Booklet-Maintenance`
- [x] 1.7 Update `src/components/MaintenancePage.tsx`: replace "Bookleaf App" with "Booklet"
- [x] 1.8 Update `src/components/SimplePageLayout.tsx`: replace "Bookleaf" nav link text with "Booklet"

## 2. New components and routing

- [x] 2.1 Create `src/pages/AppPage.tsx` — empty placeholder component
- [x] 2.2 Create `src/pages/HomePage.tsx` — login button that calls Kinde login; redirect to `/app` when already authenticated
- [x] 2.3 Create `src/components/AppLayout.tsx` — consumes `useMaintenanceActive()`; renders `<MaintenancePage />` when true, `<Outlet />` otherwise
- [x] 2.4 Create `src/App.tsx` — full route tree: `/` and `/callback` under `PublicThemeLock`, `/app` under `AuthGuard` → `AppLayout`

## 3. Unit tests

- [x] 3.1 Write `src/pages/HomePage.test.tsx`: login button visible when unauthenticated; redirects to `/app` when authenticated
- [x] 3.2 Write `src/components/AppLayout.test.tsx`: renders `MaintenancePage` when maintenance active; renders outlet content when inactive

## 4. Build and lint

- [x] 4.1 Run `npm run build` from `frontend/` and fix any type errors
- [x] 4.2 Run `npm run lint` from `frontend/` and fix any issues
