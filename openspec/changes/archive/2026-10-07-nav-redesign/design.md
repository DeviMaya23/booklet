## Context

`AppSidebar.tsx` is currently a thin wrapper around shadcn sidebar primitives with a flat `navItems` array, no icons, and no grouping. The sidebar uses `collapsible="icon"` with a `SidebarRail`. The redesign replaces the flat list with a structured layout: a standalone Inbox card at the top and two labeled nav groups below.

## Goals / Non-Goals

**Goals:**
- Introduce grouped navigation with muted uppercase group headings (LIBRARY, MANAGE)
- Add Lucide icons (18px, `strokeWidth={1.75}`) to all nav rows
- Replace the Inbox nav item with a 56px fixed-height card (icon tile, title, subtitle, count badge)
- Live Inbox count via `useFiles()` — badge hidden at zero, subtitle reads "All sorted"
- Active row: deeper-tone rounded fill, `font-semibold`
- Remove sidebar collapse (rail + `collapsible` prop)

**Non-Goals:**
- Dashboard row (route doesn't exist yet — skip entirely)
- Any backend or API changes
- Inbox page behavior changes

## Decisions

### Reuse `useFiles()` for the count
`useFiles()` fetches unassigned files (`/files?unassigned=true`) and is already called by the Inbox page. React Query deduplicates inflight requests by query key (`FILES_QUERY_KEY`), so calling it from the sidebar adds no extra network traffic when the user is on the Inbox page. On other pages it will fetch on mount and keep the count live via the existing `refetchInterval` logic.

**Alternative considered:** A separate lightweight endpoint returning only the count. Rejected — adds backend work and a new query key, while the existing hook already returns everything needed.

### Custom Inbox card, not a SidebarMenuButton
The Inbox card has a fixed 56px height, an icon tile (rounded square background), a two-line text block, and a right-edge badge. This doesn't map to `SidebarMenuButton` without fighting its layout. It will be written as a plain `NavLink`-wrapped `div` styled with Tailwind, sitting above the first `SidebarGroup`.

### Keep SidebarGroup / SidebarMenu for nav rows
The regular nav rows (Artpieces, Characters, Commissions, Artists) still fit `SidebarMenuItem` + `SidebarMenuButton`. Group headings use `SidebarGroupLabel`. This keeps shadcn's focus/keyboard semantics.

### Remove collapsible mode
The design shows no collapse control. `collapsible="icon"` and `<SidebarRail />` will be removed. `AppShell.tsx` passes no `collapsible` prop (defaults to `"offcanvas"` which is fine, or we set `collapsible="none"`).

## Risks / Trade-offs

- **`useFiles()` always active** — The file list query now runs on every authenticated page, not just Inbox. This is intentional (live count), but means one background query is always in flight. The existing `refetchInterval` only polls when thumbnails are pending, so it's not noisy.
- **shadcn sidebar primitives partially bypassed** — The Inbox card is outside the `SidebarMenu` tree. If shadcn ever adds sidebar-level keyboard nav, the card won't participate automatically. Acceptable for now.
