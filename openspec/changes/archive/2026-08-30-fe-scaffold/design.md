## Context

The frontend folder was seeded by copying files from bookleaf (the sister app). The infrastructure is largely in place: Kinde auth wiring, AuthGuard, CallbackPage, theme system, api fetch wrapper, and maintenance mode store. What's missing is `App.tsx` (the route tree entry point), the two page components, and a pass to strip all bookleaf-specific identifiers and rename them to booklet.

## Goals / Non-Goals

**Goals:**
- Wire up a minimal route tree with a public home page and a protected app shell
- Replace all bookleaf identifiers (storage keys, header names, UI strings, log prefixes) with booklet equivalents
- Integrate maintenance mode into the authenticated route layer

**Non-Goals:**
- Any real app content on `AppPage` — it is an empty placeholder
- UI design or branding beyond functional correctness (no styling work)
- Dark mode or multi-theme support — `warm` is the only theme for now

## Decisions

### Route tree structure

```
/            →  PublicThemeLock  →  HomePage
/callback    →  PublicThemeLock  →  CallbackPage
/app         →  AuthGuard        →  AppLayout  →  AppPage
```

`PublicThemeLock` wraps public routes and locks `data-theme="warm"` for the unauthenticated experience. `AuthGuard` redirects unauthenticated users to `/`. `AppLayout` sits between `AuthGuard` and the actual pages and owns the maintenance gate.

**Alternative considered:** Put the maintenance check at the root level so it covers public routes too. Rejected — the maintenance page must not block the `/callback` route. If a user is mid-auth when maintenance activates, blocking the callback traps them in an unrecoverable redirect loop.

### Maintenance mode placement

`AppLayout` consumes `useMaintenanceActive()` and short-circuits to `<MaintenancePage />` when true. All future routes added under `/app` get this automatically, with no per-page wiring required.

`apiFetch` (already in `api.ts`) reads the `X-Booklet-Maintenance` response header on every call and calls `setMaintenanceActive(true)`. Since `AppPage` is currently empty and makes no API calls, the maintenance gate will be inert until the first real endpoint is wired up — that is expected and correct.

### Theme system

Keep `ThemeProvider` and `useTheme` as-is structurally, but narrow the `Theme` type to `'warm'` only and change the storage key to `booklet-theme`. The `PublicThemeLock` continues to hardcode `data-theme="warm"` for public routes. The in-HTML theme-restore script in `index.html` must also reference `booklet-theme` so the pre-render theme read is consistent.

**Alternative considered:** Drop the theme system entirely for now and hardcode. Rejected — `ThemeProvider` wraps the entire tree in `main.tsx`; removing it requires more surgery than updating two string values.

### bookleaf → booklet identifier sweep

Every bookleaf-specific identifier in the copied files is renamed to its booklet equivalent:

| Old | New |
|-----|-----|
| `bookleaf-theme` (localStorage) | `booklet-theme` |
| `bookleaf-maintenance-bypass` (localStorage) | `booklet-maintenance-bypass` |
| `X-Bookleaf-Bypass` (header) | `X-Booklet-Bypass` |
| `X-Bookleaf-Maintenance` (header) | `X-Booklet-Maintenance` |
| `[bookleaf]` (console prefix) | `[booklet]` |
| `"Bookleaf"` / `"Bookleaf App"` (UI strings) | `"Booklet"` / `"Booklet"` |
| `<title>Bookleaf</title>` | `<title>Booklet</title>` |

## Risks / Trade-offs

- **Maintenance mode is inert at launch** — `apiFetch` is the trigger, but `AppPage` makes no API calls yet. The gate is correctly wired; it just won't fire until the first real feature is added. This is acceptable and expected.
- **Single theme** — narrowing `Theme` to `'warm'` means any stored `'lumen'` or `'sunless'` values from a bookleaf session on the same origin would fall back to `'warm'`. No real users exist yet, so there is no migration concern.
