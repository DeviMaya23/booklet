## 1. Backend — Dashboard Repository

- [x] 1.1 Create `DashboardRepository` interface in `backend/internal/usecase/` with methods: `GetRecentArtpieces`, `GetInProgressCommissions`, `GetCommissionsNoArtist`, `GetArtpiecesNoArtist`, `GetDoneCommissionsNoArtpieces`, `GetArtpiecesNoFiles`
- [x] 1.2 Implement `GetRecentArtpieces` in `backend/internal/repository/dashboard_repository.go` — query artpieces ordered by `created_at DESC LIMIT 5`, preload `Artist` and `CoverFile`
- [x] 1.3 Implement `GetInProgressCommissions` — query commissions where `status IN ('waitlist', 'wip')` ordered by `created_at ASC`, preload `Artist`
- [x] 1.4 Implement `GetCommissionsNoArtist` — query commissions where `artist_id IS NULL`, return id + title only
- [x] 1.5 Implement `GetArtpiecesNoArtist` — query artpieces where `artist_id IS NULL`, return id + title only
- [x] 1.6 Implement `GetDoneCommissionsNoArtpieces` — query commissions where `status = 'done'` and `NOT EXISTS (SELECT 1 FROM artpieces WHERE commission_id = commissions.id)`, return id + title only
- [x] 1.7 Implement `GetArtpiecesNoFiles` — query artpieces where `NOT EXISTS (SELECT 1 FROM files WHERE artpiece_id = artpieces.id)`, return id + title only

## 2. Backend — Dashboard Usecase and Handler

- [x] 2.1 Create `DashboardUsecase` in `backend/internal/usecase/dashboard_usecase.go` with a `GetDashboard(ctx, userID)` method that calls all six repo methods sequentially and assembles the response struct
- [x] 2.2 Define response structs: `DashboardResponse`, `DashboardRecentArtpiece`, `DashboardInProgressCommission`, `DashboardHousekeeping`, and housekeeping item struct
- [x] 2.3 Create `DashboardHandler` in `backend/internal/handler/dashboard_handler.go` with `GetDashboard` method — authenticate user, call usecase, generate presigned thumbnail URLs for recent artpieces, return JSON
- [x] 2.4 Register `GET /dashboard` route in `backend/cmd/server/main.go`, wire up handler and repository dependencies
- [x] 2.5 Create Bruno file `backend/bruno/dashboard/get-dashboard.bru`

## 3. Backend — Unit Tests

- [x] 3.1 Write unit tests for `DashboardUsecase.GetDashboard` — test successful aggregation of all six queries, and test each repo method error propagating correctly

## 4. Backend — Lint

- [x] 4.1 Run `golangci-lint run ./...` from the backend directory and fix any issues

## 5. Frontend — Dashboard API Hook

- [x] 5.1 Create `frontend/src/features/dashboard/api/useDashboard.ts` — React Query hook fetching `GET /dashboard`, query key `['dashboard']`
- [x] 5.2 Define TypeScript types matching the dashboard response shape: `DashboardData`, `DashboardArtpiece`, `DashboardCommission`, `DashboardHousekeeping`

## 6. Frontend — Dashboard Page

- [x] 6.1 Create `frontend/src/pages/DashboardPage.tsx` — top-level page component, mounts `useDashboard()`, passes data to section components
- [x] 6.2 Implement header row: "Upload" button (opens `AddFilesModal`) and "+ New Commission" button (no-op)
- [x] 6.3 Implement Recent Artpieces section: horizontal thumbnail row using existing artpiece thumbnail style, "View all" link to `/app/artpieces`; artpiece tile shows thumbnail (placeholder if null), title, artist name
- [x] 6.4 Implement In Progress panel: commission rows showing title, artist name (with artist-link affordance), status chip, paid state, days-elapsed (computed from `created_at`), last-contacted timestamp; stale highlighting (orange) when `last_contacted_at` > 14 days ago; muted color when null; "View all" link to `/app/commissions`
- [x] 6.5 Implement mark-contacted action on each In Progress row: calls `usePatchCommission` directly with `lastContactedAt: new Date().toISOString()`; after the mutation resolves, call `queryClient.invalidateQueries(['dashboard'])` locally in the dashboard component — do not add this invalidation to `usePatchCommission` itself
- [x] 6.6 Implement Housekeeping panel: four category sections each with label and list of linked items (commission → `/app/commissions`, artpiece → `/app/artpieces/:id`); all-clear empty state ("All tidy. Nothing needs fixing.") when all four arrays are empty
- [x] 6.7 Add loading state (skeleton or spinner) while dashboard data is fetching

## 7. Frontend — Routing and Sidebar

- [x] 7.1 Add `/app/dashboard` route in `App.tsx`, update `/app` redirect from `/app/characters` to `/app/dashboard`
- [x] 7.2 Add Dashboard nav item to `AppSidebar.tsx` above the Inbox card, using `layout-dashboard` Lucide icon; apply same active-state styling as other nav rows

## 8. Frontend — Build and Lint

- [x] 8.1 Run `npm run build` from the frontend directory and fix any issues
- [x] 8.2 Run `npm run lint` from the frontend directory and fix any issues
