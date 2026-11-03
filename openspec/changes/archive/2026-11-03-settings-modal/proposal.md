# Proposal

## Why

Clicking the avatar in the top bar immediately logs the user out, which is easy to trigger by accident and leaves no home for account-level actions. The backend already supports account deletion (`DELETE /me`) and the CSS already defines three themes (`warm`, `lumen`, `sunless`), but neither is reachable from the UI. A profile dropdown plus a settings modal (modelled on Bookleaf's) wires up both.

## What Changes

- Clicking the top-bar avatar opens a dropdown (email label, **Settings**, separator, red **Sign out**) instead of logging out.
- New Settings modal opened from the dropdown, with a left nav of two sections: **Account** and **App**. Advanced and Extensions are intentionally omitted.
- **Account** section: profile block (avatar, full name, read-only email), a Log out button, and a Delete account flow with an inline type-your-email confirmation. A successful `DELETE /me` is followed by Kinde `logout()` (so the user does not hit the "Please log back in" session-expired toast).
- **App** section: Theme picker with three radio cards (Parchment = `warm`, Lumen, Sunless), each with swatches and a subtitle.
- `useTheme` is widened from `'warm'` only to `'warm' | 'lumen' | 'sunless'`; stored values are validated against that set.
- The Settings dropdown item is hidden on mobile viewports (<768px), using the existing `useIsMobile` hook.
- The existing avatar markup (image, or initial fallback) is reused in the top bar and the Account profile block; no new `ui/avatar` component is added.
- **BREAKING (spec only)**: `web-app-shell` currently says clicking the avatar logs out; that requirement is replaced.

No backend changes: `DELETE /me` already exists (spec `account-deletion`).

## Capabilities

### New Capabilities
- `web-settings-modal`: the profile dropdown's Settings entry, the modal shell and section nav (Account, App), the Account section (profile, log out, delete account with email confirmation), and the App section's theme picker UI.
- `web-theme`: theme preference model: the three selectable themes, persistence across reloads, fallback for invalid stored values, and application to the authenticated app.

### Modified Capabilities
- `web-app-shell`: the "Top bar" requirement changes from "avatar click logs out" to "avatar click opens the profile dropdown", including the Sign out item and the Settings item.

## Impact

- Frontend only (`frontend/src`):
  - `components/AppTopBar.tsx`: replace logout button with dropdown trigger; reuse existing avatar markup.
  - New `features/settings/` (modal, Account and App sections, delete-account mutation hook following the `features/*/api` pattern).
  - `hooks/useTheme.tsx` and `hooks/useTheme.test.tsx`: widen `Theme`, validate stored value.
- Existing `components/ui` primitives only (`dropdown-menu`, `dialog`, `button`, `input`); no new dependencies or `ui` components.
- Call sites of the shared symbols being changed: `useTheme`/`Theme` is only consumed by `ThemeProvider` in `main.tsx` today; the `logout` call in `AppTopBar` is the only avatar logout entry point.
- API: consumes existing `DELETE /me`. No new endpoint, so no Bruno file.
- Docs/specs: delta specs for `web-app-shell`, `web-settings-modal`, `web-theme`.
