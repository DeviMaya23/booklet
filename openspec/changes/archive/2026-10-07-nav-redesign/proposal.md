## Why

The current sidebar is a plain flat list with no icons, no grouping, and no contextual information. The redesign introduces visual hierarchy, grouped navigation, and a live Inbox count card that surfaces actionable state (unassigned files) at a glance.

## What Changes

- Replace the flat nav list with two labeled groups: **Library** (Artpieces, Characters) and **Manage** (Commissions, Artists)
- Add Lucide icons (18px, stroke 1.75) to every nav row
- Replace the bare "Inbox" nav item with a 56px card that shows a live count of unassigned files and a contextual subtitle ("X to sort" / "All sorted")
- Remove the collapsible sidebar rail — the sidebar is always expanded
- Active row style: deeper-tone rounded fill, heavier font weight
- Group headings: small muted uppercase labels, no divider lines

## Capabilities

### New Capabilities
- `web-nav-sidebar`: Grouped, icon-labelled sidebar navigation with an Inbox count card

### Modified Capabilities
- `web-app-shell`: Sidebar navigation requirements change — groups, icons, Inbox card, no collapse

## Impact

- `frontend/src/components/AppSidebar.tsx` — full rewrite
- `frontend/src/components/AppShell.tsx` — remove `SidebarRail`, possibly remove `collapsible="icon"` prop
- Reuses `useFiles()` hook from `web-file-inbox` — no new API calls
- No backend changes
