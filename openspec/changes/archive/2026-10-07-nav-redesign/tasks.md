## 1. AppShell cleanup

- [x] 1.1 Remove `collapsible="icon"` from `SidebarProvider` in `AppShell.tsx` (or set `collapsible="none"`)
- [x] 1.2 Remove `<SidebarRail />` from `AppSidebar.tsx`

## 2. Inbox card

- [x] 2.1 Call `useFiles()` inside `AppSidebar` and derive the unassigned count from `data?.length ?? 0`
- [x] 2.2 Build the `InboxCard` component: 56px fixed height, icon tile, bold "Inbox" title, subtitle ("X to sort" / "All sorted"), dark count badge hidden at zero
- [x] 2.3 Wire `InboxCard` as a `NavLink` to `/app/files` and place it above the nav groups

## 3. Nav groups

- [x] 3.1 Replace the flat `navItems` array with two group definitions: LIBRARY (Artpieces, Characters) and MANAGE (Commissions, Artists), each entry including its Lucide icon component
- [x] 3.2 Render each group using `SidebarGroup` + `SidebarGroupLabel` (muted uppercase) + `SidebarMenu`
- [x] 3.3 Add icon to each `NavItem` — render the passed icon at 18px with `strokeWidth={1.75}`
- [x] 3.4 Apply active row style: deeper-tone rounded fill + `font-semibold` when route matches

## 4. Verification

- [x] 4.1 Run `npm run build` and fix any type errors
- [x] 4.2 Run `npm run lint` and fix any lint issues
