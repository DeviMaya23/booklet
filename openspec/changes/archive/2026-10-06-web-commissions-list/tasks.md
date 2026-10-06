## 1. Backend — migration and domain

- [x] 1.1 Create migration `000027_add_last_contacted_to_commissions` — `ALTER TABLE commissions ADD COLUMN last_contacted_at TIMESTAMPTZ`
- [x] 1.2 Add `LastContactedAt *time.Time` to the `Commission` domain struct in `domain/commission.go`

## 2. Backend — commission response update

- [x] 2.1 Add `ArtistLink *string` and `LastContactedAt *string` fields to `commissionResponse` in `commission_handler.go`
- [x] 2.2 Map both fields in `toCommissionResponse`: `ArtistLink` from `c.Artist.ArtistLink` (guard for nil artist), `LastContactedAt` formatted as RFC3339

## 3. Backend — PATCH endpoint

- [x] 3.1 Add `PatchCommissionParams` to `commission_repository.go` with pointer fields: `Status *string`, `Paid *bool`, `PaidDate *string`, `LastContactedAt *string`; add `Patch` method to `CommissionRepository` interface
- [x] 3.2 Implement `Patch` in `repository/commission_repository.go` — build the update map from non-nil fields only and call `db.Model(...).Updates(map)`
- [x] 3.3 Add `Patch` method to `CommissionUsecase` and `CommissionUsecase` interface in handler — validate `Status` against allowed enum if non-nil; return 422 for invalid value
- [x] 3.4 Add `patchCommissionRequest` struct and `PatchCommission` handler method in `commission_handler.go` — use `c.Bind` + `c.Validate`; map to `PatchCommissionParams`; return updated commission with HTTP 200
- [x] 3.5 Register `PATCH /commissions/:id` route in `main.go`
- [x] 3.6 Add bruno file `collection/commissions/patch_commission.bru`

## 4. Backend — unit tests

- [x] 4.1 Unit test `PatchCommission` usecase: status validation (invalid value returns `ErrInvalidCommissionStatus`); commission not found returns error; successful patch with each individual field

## 5. Frontend — shadcn components

- [x] 5.1 Install shadcn `table` component (`npx shadcn@latest add table`)
- [x] 5.2 Install shadcn `checkbox` component (`npx shadcn@latest add checkbox`)
- [x] 5.3 Install shadcn `calendar` and `popover` components (`npx shadcn@latest add calendar popover`)

## 6. Frontend — commissions API hooks

- [x] 6.1 Create `features/commissions/api/useCommissions.ts` — `GET /commissions` query; define `Commission` type including `artist_link`, `last_contacted_at`; export `COMMISSIONS_QUERY_KEY`
- [x] 6.2 Create `features/commissions/api/usePatchCommission.ts` — `PATCH /commissions/:id` mutation accepting `{ id, status?, paid?, paidDate?, lastContactedAt? }`; invalidates `COMMISSIONS_QUERY_KEY` on success

## 7. Frontend — utility

- [x] 7.1 Create `features/commissions/utils/relativeTime.ts` — pure function `formatRelativeTime(date: Date): string` returning "today", "yesterday", "N days ago", "N weeks ago", "N months ago"

## 8. Frontend — CommissionsTable component

- [x] 8.1 Create `features/commissions/components/CommissionsTable.tsx` — accepts `commissions: Commission[]` and `view: 'active' | 'done'`; renders shadcn `Table` with columns: Title, Artist (name + link icon), Status, Paid, Paid Date, and the view-dependent column; all columns sortable except artist link
- [x] 8.2 Implement Status cell — shadcn `DropdownMenu` inline; on change call `usePatchCommission`; show inline error border + icon on mutation error
- [x] 8.3 Implement Paid cell — shadcn `Checkbox`; on check-to-true send `{ paid: true, lastContactedAt: new Date().toISOString() }`; on uncheck send `{ paid: false }`; show inline error on failure
- [x] 8.4 Implement Paid Date cell — shadcn `Popover` + `Calendar` datepicker; on date pick call `usePatchCommission` with selected date; show inline error on failure
- [x] 8.5 Implement Last Contact cell (active view) — display `formatRelativeTime(last_contacted_at)` or em dash if null; stamp button (circular icon) calls `usePatchCommission` with `lastContactedAt: new Date().toISOString()`; show inline error on failure
- [x] 8.6 Implement Time Taken cell (done view) — compute `finish_date - paid_date` in days; display "N days" or em dash if either is null

## 9. Frontend — CommissionsPage and routing

- [x] 9.1 Create `pages/CommissionsPage.tsx` — toggle state (`active` | `done`, default `active`); filter `commissions` client-side by status; render page header with toggle and "+ New Commission" button (no-op); render `CommissionsTable`
- [x] 9.2 Add `/app/commissions` route in `App.tsx`
- [x] 9.3 Add "Commissions" entry to the `navItems` array in `AppSidebar.tsx`

## 10. Final checks

- [x] 10.1 Run `golangci-lint run` in the backend and fix any issues
- [x] 10.2 Run `npm run build` in the frontend and fix any type errors
- [x] 10.3 Run `npm run lint` in the frontend and fix any issues
