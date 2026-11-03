# Design

## Context

See proposal.md for motivation. Current state, all verified in the repo:

- `AppTopBar` renders an avatar `<button>` that calls Kinde `logout()` directly, using `user` from `useKindeAuth()` and hand-rolled `<img>` / initial-fallback markup.
- `DELETE /me` is complete on the backend (`UserHandler.DeleteMe` -> `MarkPendingDeletion` -> Bookleaf, returns 202). After it succeeds the account is `pending_deletion`, so any further API call returns 401, which `apiFetch` turns into the "Please log back in" toast via `sessionExpiredStore`. Nothing in the FE calls it yet.
- `index.css` already defines `[data-theme="warm" | "lumen" | "sunless"]` token sets, and the `dark:` variant is keyed to `sunless`. `ThemeProvider` writes `data-theme` on `<html>` and persists to localStorage key `booklet-theme`, but `Theme` is typed `'warm'` and `readStoredTheme` accepts only `'warm'`. Nothing calls `setTheme`.
- `components/ui` already has `dropdown-menu`, `dialog`, `button`, `input` (Base UI). `DialogContent` here has no `showCloseButton` prop (Bookleaf's does) and renders no built-in close button.
- Bookleaf's `ProfileMenu`, `SettingsModal`, `AccountSection` and `AppSection` are the visual and behavioural reference. Booklet differs: no `ui/avatar`, no `PATCH /me` (so no folder-icons toggle), no Advanced/Extensions.

## Goals / Non-Goals

**Goals:**
- Replace avatar-click-logout with a dropdown, add the Settings modal with Account and App sections, match Bookleaf's look.
- Make the three existing themes selectable and persisted.
- Wire the existing `DELETE /me` with a safe confirmation.

**Non-Goals:**
- No backend changes, no new endpoints, no Bruno files.
- No Advanced or Extensions sections, no folder-icons toggle, no server-side theme storage (stays per-browser localStorage, as today).
- No new `ui` components or dependencies (no `avatar`, `radio-group`, `tabs`).
- No changes to Bookleaf.
- No changes to public pages or `PublicThemeLock` (the login page is still a stub).

## Decisions

**1. Feature folder `features/settings/`.** Holds `components/SettingsModal.tsx`, `AccountSection.tsx`, `AppSection.tsx`, `api/useDeleteAccount.ts`, and a confirmation hook. `AppTopBar` stays in `components/` and owns the `settingsOpen` state, as Bookleaf's `ProfileMenu` does. Alternative: put it under `features/auth`. Rejected, since auth currently only holds the guard and profile helpers, and settings is a distinct feature surface.

**2. Delete-account mutation follows the existing `features/*/api` hook pattern** (`useDeleteAccount`, using `useKindeAuth().getToken` + `apiFetch('/me', getToken, { method: 'DELETE' })`, throwing on non-ok), like `useDeleteCharacter`. The caller does the Kinde `logout()` in `onSuccess` and the "Failed to delete account" toast in `onError`.

**3. Confirmation state is its own hook** (`useDeleteAccountConfirmation(email)`), per the one-hook-per-concern convention. It returns `confirming`, `typedEmail`, `emailMatches`, and semantic mutators `start()`, `cancel()`, `updateTypedEmail(v)`. Matching is `trim().toLowerCase()` equality; it is false when the email is unknown.

**4. Logout after delete, not a redirect.** On 202 call `logout()` immediately so the user never makes a request as a `pending_deletion` account (which would 401 into the session-expired toast). Keep the dialog/prompt in place on error.

**5. Profile data comes from `useKindeAuth().user`**, already used by `AppTopBar`, not Bookleaf's `getUserProfile` query. This is synchronous and avoids a new query key. `profile.ts`'s `getInitials`/`getFullName` take a `UserProfile`; the Kinde `user` object has the same name fields, so they are reused if the types line up, otherwise a trivial adapter inside `features/settings` is used (checked at implementation).

**6. Avatar markup is reused inline, not extracted.** The `<img>` / initial fallback from `AppTopBar` is repeated in the Account profile block with a larger size. Two occurrences stay within the conventions' "don't abstract on the first or second occurrence" rule, and the user chose to reuse the current markup over adding `ui/avatar`.

**7. Mobile hiding via the existing `useIsMobile` hook** (Booklet's own 768px breakpoint, also used by the sidebar) conditionally not rendering the Settings item. Bookleaf used a CSS `hidden sm:flex` (640px); the hook is already in the codebase and makes the behaviour testable in jsdom. Trade-off: breakpoint is 768px rather than Bookleaf's 640px. If the modal is open while the viewport shrinks, it simply stays open (no forced close).

**8. Theme picker uses native radio inputs inside labels**, as Bookleaf does. Labels, subtitles and swatch hex values live in a small constant map in `AppSection`, keyed by the `Theme` union. The stored key stays `warm`; "Parchment" is display-only.

**9. Widen `Theme` and validate stored values against a const list** (`['warm','lumen','sunless']`) rather than the current `=== 'warm'` check, so an invalid or legacy value still falls back to `warm`.

## Risks / Trade-offs

- [`useKindeAuth().user` lacks `email` or the fields `getInitials` expects] -> Check the Kinde type at implementation; fall back to `getUserProfile()` via react-query as Bookleaf does.
- [Dialog content pane styling differs from Bookleaf because Booklet's `DialogContent` has no `showCloseButton` prop] -> Omit the prop; Booklet's dialog has no built-in close button, so the custom `DialogClose` is the only one.
- [`sunless` was never selectable before, so some components may have unverified dark-theme styling] -> Browse each main page under each theme during the manual check; fix only token-level issues, and flag anything larger rather than expanding scope.
- [Theme remains per-browser] -> Accepted, unchanged from today.

## Open Questions

None blocking.
