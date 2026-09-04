## 1. Initialize shadcn

- [x] 1.1 Run `npx shadcn@latest init` inside `frontend/` — when prompted for style, select the `base-` variant (Base UI); set component output path to `src/components/ui`
- [x] 1.2 Review generated `components.json` and confirm style is `base-*`
- [x] 1.3 Review any changes made to `index.css` — ensure the four theme blocks (`warm`, `lumen`, `sunless`, and the dark custom variant) are intact; restore if overwritten

## 2. Install sidebar block

- [x] 2.1 Run `npx shadcn@latest add sidebar` inside `frontend/`
- [x] 2.2 Confirm `src/components/ui/sidebar.tsx` was generated

## 3. Align sidebar CSS tokens

- [x] 3.1 In `index.css`, update `--sidebar-accent`, `--sidebar-border`, and `--sidebar-ring` in the `:root` / `[data-theme="warm"]` blocks to use the warm palette (`rgba(45, 42, 38, ...)` family, consistent with `--border` and `--accent`)
- [x] 3.2 Apply the same corrections to the `[data-theme="lumen"]` block
- [x] 3.3 Apply the same corrections to the `[data-theme="sunless"]` block

## 4. Page stubs

- [x] 4.1 Create `src/pages/CharactersPage.tsx` — `export default function CharactersPage() { return null }`
- [x] 4.2 Create `src/pages/ImagesPage.tsx` — same pattern
- [x] 4.3 Create `src/pages/ArtistsPage.tsx` — same pattern

## 5. AppSidebar

- [x] 5.1 Create `src/components/AppSidebar.tsx` — use the installed shadcn sidebar components to render a collapsible sidebar with three nav items: Characters (`/app/characters`), Images (`/app/images`), Artists (`/app/artists`); use `NavLink` from `react-router-dom` for active state

## 6. AppTopBar

- [x] 6.1 Create `src/components/AppTopBar.tsx` — render "booklet" as a serif wordmark on the left; on the right, render a user avatar using `user.picture` from `useKindeAuth()`; clicking the avatar calls `logout()`

## 7. AppShell

- [x] 7.1 Create `src/components/AppShell.tsx` — compose `AppTopBar` and `AppSidebar`; render `<Outlet />` as the main content area

## 8. Route wiring

- [x] 8.1 In `src/App.tsx`, replace `<Route path="/app" element={<AppPage />} />` with a nested structure under `AppShell`: index route redirects to `/app/characters`; add child routes for `/app/characters`, `/app/images`, `/app/artists`
- [x] 8.2 Delete `src/pages/AppPage.tsx`

## 9. Build and lint

- [x] 9.1 Run `npm run build` from `frontend/` and fix any type errors or build failures
- [x] 9.2 Run `npm run lint` from `frontend/` and fix any reported issues
