# Tasks

## 1. Theme preference (`web-theme`)

- [x] 1.1 Widen `Theme` in `frontend/src/hooks/useTheme.tsx` to `'warm' | 'lumen' | 'sunless'` and validate the stored value against a const list (invalid -> `warm`). Update `useTheme.test.tsx`: restores a stored `lumen`/`sunless`, `setTheme('sunless')` updates `data-theme` and localStorage, an unrecognised stored value falls back to `warm`. Verify `npx vitest run src/hooks/useTheme.test.tsx` passes.

## 2. Delete-account logic (`web-settings-modal`)

- [x] 2.1 Add `frontend/src/features/settings/api/useDeleteAccount.ts` (`DELETE /me` via `apiFetch`, throw on non-ok, same shape as `useDeleteCharacter`) with a colocated test that mocks the API boundary inside a fresh `QueryClientProvider` (`retry: false`): 202 resolves, 500 rejects. Verify the test passes.
- [x] 2.2 Add `useDeleteAccountConfirmation(email)` in `features/settings/` returning `confirming`, `typedEmail`, `emailMatches` and semantic mutators `start`, `cancel`, `updateTypedEmail`. Colocated test covers: non-matching text -> false; matching ignoring case and surrounding whitespace -> true; unknown email -> false; `cancel` clears text and exits confirming. Verify the test passes.

## 3. Settings modal UI (`web-settings-modal`)

- [x] 3.1 Build `AccountSection` (profile block reusing the existing avatar image/initial-fallback markup, full name, read-only email, Log out button, Delete account area using tasks 2.1 and 2.2). On success call Kinde `logout()` and not the session-expired flow; on error show `toast.error('Failed to delete account')` and keep the prompt. Component tests (mock `useKindeAuth`, `sonner` and the API boundary): profile name and email shown; Log out calls `logout`; prompt reveal; confirm disabled until email matches; Cancel resets; success calls `logout`; failure shows the toast, does not call `logout`, and keeps the prompt visible. Verify tests pass.
- [x] 3.2 Build `AppSection` with the Parchment / Lumen / Sunless radio cards (swatches, label, subtitle; preselect current theme; selecting calls `setTheme`). Component tests inside `ThemeProvider`: current theme preselected, all three listed with subtitles, selecting Lumen changes the selection and `document.documentElement.dataset.theme`. Verify tests pass.
- [x] 3.3 Build `SettingsModal` (`Dialog`, left nav with exactly Account and App, section heading, `DialogClose` button; opens on Account; no Advanced/Extensions; no `showCloseButton` prop). Component tests: opens on Account, switching to App shows its heading, only two nav entries, close button calls `onOpenChange(false)`. Verify tests pass.

## 4. Top bar dropdown (`web-app-shell`, `web-settings-modal`)

- [x] 4.1 Rework `frontend/src/components/AppTopBar.tsx`: avatar becomes a `DropdownMenu` trigger (keep an accessible name), content shows email label, Settings (omitted when `useIsMobile()` is true), separator, red Sign out (`logout()`); owns `settingsOpen` and renders `SettingsModal`. Add `AppTopBar.test.tsx`: clicking the avatar opens the dropdown and does not call `logout`; Sign out calls `logout`; email shown; Settings opens the modal on Account; Settings absent on a mobile viewport and present on desktop; initial fallback shown with no picture and `?` with no name. Verify tests pass.
- [x] 4.2 Grep for every other consumer of the avatar logout button, `Log out` aria-label and `useTheme`/`Theme` across `frontend/src` (including tests) and confirm nothing else depended on the old behaviour; fix any stale test or reference found.

## 5. Integration checks

- [x] 5.1 Run `npm run build` and `npm run lint` in `frontend/` and fix whatever arises (no `eslint-disable` suppressions).
- [x] 5.2 Manual check in the browser: avatar opens dropdown, Settings opens the modal, theme switches live and survives reload under all three themes, signing out returns to `/`, Settings is hidden in a narrow window, and the delete flow's confirm button stays disabled until the email matches (do not execute a real deletion against a live account).
- [x] 5.3 Run `openspec validate settings-modal --strict` and confirm it passes.

## Workflow follow-up

- Archive the change after review so the `web-app-shell` delta is merged and `web-settings-modal` and `web-theme` become main specs.
