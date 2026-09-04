## Context

The frontend has an authenticated route at `/app` guarded by `AuthGuard` and wrapped by `AppLayout` (which handles the maintenance gate). `AppPage` — the current leaf — returns null. This change builds the shell layer that will persist across all authenticated views and establishes the routing structure for the three main sections.

The stack uses `@base-ui/react` for UI primitives (not Radix). shadcn is installed as a CLI dependency (`"shadcn": "^4.7.0"`) but has not been initialized; no `components.json` exists yet. The CSS already contains `--sidebar-*` tokens copied from a reference project, but some values use generic oklch defaults rather than the project's warm palette.

## Goals / Non-Goals

**Goals:**
- Deliver a working app shell at `/app/*` with a persistent top bar and collapsible sidebar
- Wire three stub routes (characters, images, artists) so the nav is functional end to end
- Use shadcn's Base UI sidebar block as the implementation foundation

**Non-Goals:**
- Any page content (characters list, image grid, artist list) — stubs only
- Mobile/responsive breakpoints beyond what the shadcn block ships with
- Theme switcher in the top bar
- User profile or settings page

## Decisions

### shadcn init with Base UI style

`shadcn init` is run with the `base-` style so generated components use `@base-ui/react` primitives. This satisfies the CONVENTIONS.md constraint (no Radix, no `asChild`).

_Alternative considered_: Hand-roll the sidebar from scratch using `@base-ui/react` primitives directly. Rejected — the shadcn Base UI sidebar block gives collapse behavior, keyboard navigation, and accessible markup that would take significant effort to replicate correctly.

### AppShell as a layout route

`AppShell` is a layout component that renders inside `AppLayout` and outputs the top bar, sidebar, and `<Outlet />`. `AppPage` is removed.

```
AuthGuard
  AppLayout        ← maintenance gate (unchanged)
    AppShell       ← new: top bar + sidebar + <Outlet />
      CharactersPage | ImagesPage | ArtistsPage
```

_Alternative considered_: Keep `AppPage` and render the shell inside it. Rejected — the shell is layout, not a page. Modeling it as a layout route is consistent with `AppLayout` and makes the tree readable.

### Three separate components

```
AppShell      ← composes the two; renders <Outlet />
AppTopBar     ← "booklet" wordmark + Kinde user avatar with logout
AppSidebar    ← shadcn sidebar block; NavLink items for three sections
```

Top bar and sidebar are split because they have independent concerns and will grow independently (top bar may gain global search; sidebar may gain section badges or nesting).

### Sidebar CSS token alignment

`index.css` was copied from a reference project. The `--sidebar-accent`, `--sidebar-border`, and `--sidebar-ring` values use generic oklch defaults rather than the warm palette. These are corrected in this change since the sidebar is first rendered here — leaving them misaligned would surface visual inconsistency immediately.

### No unit tests for AppShell

The shell is pure layout composition. There is no branching logic beyond a `<Navigate>` redirect and a `useKindeAuth()` call. Unit tests would only assert JSX structure, not behavior. No tests are planned.

## Risks / Trade-offs

- **shadcn init touching index.css** → `index.css` is a copy-paste starting point (not hand-crafted), so regeneration is acceptable. The four theme blocks must be verified intact after init and restored if overwritten.
- **Base UI sidebar block API unknown until installed** → The exact component exports are not confirmed until `shadcn add sidebar` runs. `AppSidebar` adapts to whatever the block ships; this is a first-install task, not a migration risk.
- **Sidebar token misalignment across themes** → Warm, lumen, and sunless theme blocks all need the corrected sidebar token values. Missing one would cause a visual regression on that theme.
