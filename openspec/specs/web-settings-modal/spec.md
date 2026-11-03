# web-settings-modal Specification

## Purpose

Defines the profile dropdown opened from the top-bar avatar and the Settings modal reached from it, including the Account section (profile, log out, delete account) and the App section's theme picker UI.

## Requirements

### Requirement: Profile dropdown
Clicking the top-bar avatar SHALL open a dropdown containing, in order: the signed-in user's email as a muted non-interactive label, a "Settings" item, a separator, and a destructive-styled "Sign out" item. The avatar SHALL show the user's picture, or the first initial of their given name when no picture exists.

#### Scenario: Dropdown shows email, Settings and Sign out
- **WHEN** an authenticated user clicks the top-bar avatar
- **THEN** the dropdown SHALL show the user's email, a "Settings" item, and a "Sign out" item, with Sign out last

#### Scenario: Sign out logs the user out
- **WHEN** the user selects "Sign out" in the dropdown
- **THEN** the app SHALL log the user out

#### Scenario: Avatar falls back to initial
- **WHEN** the user has no profile picture
- **THEN** the avatar SHALL show the uppercased first letter of their given name, or "?" when none is available

### Requirement: Settings entry is desktop only
The "Settings" dropdown item SHALL NOT be shown on mobile viewports (narrower than 768px). "Sign out" SHALL remain available at all viewport widths.

#### Scenario: Settings hidden on mobile
- **WHEN** the viewport is narrower than 768px and the user opens the dropdown
- **THEN** the dropdown SHALL NOT show a "Settings" item
- **AND** it SHALL still show the email and "Sign out"

#### Scenario: Settings shown on desktop
- **WHEN** the viewport is 768px or wider and the user opens the dropdown
- **THEN** the dropdown SHALL show the "Settings" item

### Requirement: Settings modal shell
Selecting "Settings" SHALL open a modal titled "Settings" with a left navigation containing exactly two sections, "Account" and "App", and a content pane headed with the active section's name and a close button. The modal SHALL open on the Account section. The modal SHALL NOT offer Advanced or Extensions sections.

#### Scenario: Settings opens on Account
- **WHEN** the user selects "Settings" in the dropdown
- **THEN** the Settings modal SHALL open with the "Account" section active and its heading shown

#### Scenario: Switching sections
- **WHEN** the user selects "App" in the modal's left navigation
- **THEN** the content pane SHALL show the App section and its heading SHALL read "App"

#### Scenario: Only Account and App are offered
- **WHEN** the Settings modal is open
- **THEN** the left navigation SHALL list "Account" and "App" and no other sections

#### Scenario: Closing the modal
- **WHEN** the user activates the close button
- **THEN** the Settings modal SHALL close and the user SHALL remain on the current route

### Requirement: Account profile and log out
The Account section SHALL show a profile block with the user's avatar, full name, and read-only email, followed by a "Log out" button.

#### Scenario: Profile details are shown
- **WHEN** the Account section is displayed for a user with a full name and email
- **THEN** it SHALL show that full name and email, with the email not editable

#### Scenario: Log out from Account
- **WHEN** the user clicks "Log out" in the Account section
- **THEN** the app SHALL log the user out

### Requirement: Delete account with email confirmation
The Account section SHALL include a "Delete account" area that initially shows a "Delete account…" button. Activating it SHALL reveal an inline prompt asking the user to type their account email, with a Cancel button and a "Delete account" button that is disabled until the typed text equals the user's email (trimmed, case-insensitive). Confirming SHALL call `DELETE /me`.

#### Scenario: Confirmation prompt is revealed
- **WHEN** the user clicks "Delete account…"
- **THEN** the section SHALL show the user's email in the prompt, an email input, a Cancel button, and a disabled "Delete account" button

#### Scenario: Confirm button enables only on matching email
- **WHEN** the user types text that does not match their email
- **THEN** the "Delete account" confirm button SHALL remain disabled
- **AND WHEN** the typed text matches their email ignoring case and surrounding whitespace
- **THEN** the confirm button SHALL become enabled

#### Scenario: Cancel resets the prompt
- **WHEN** the user has typed into the confirmation input and clicks Cancel
- **THEN** the prompt SHALL be hidden, the input cleared, and the "Delete account…" button shown again

#### Scenario: Successful deletion logs the user out
- **WHEN** the user confirms and `DELETE /me` responds 202
- **THEN** the app SHALL log the user out
- **AND** SHALL NOT show the "Please log back in" session-expired toast

#### Scenario: Deletion fails
- **WHEN** the user confirms and `DELETE /me` responds with a non-2xx status
- **THEN** the app SHALL show an error toast "Failed to delete account"
- **AND** the user SHALL remain signed in with the confirmation prompt still shown

### Requirement: App section theme picker
The App section SHALL show a "Theme" group with one radio card per theme: "Parchment" ("Warm parchment"), "Lumen" ("Bright and clean"), and "Sunless" ("Mostly black"). Each card SHALL show colour swatches, its label and its subtitle. The card for the current theme SHALL be selected. Selecting a card SHALL change the active theme immediately as defined in the `web-theme` spec.

#### Scenario: Current theme is preselected
- **WHEN** the App section is displayed while the active theme is Sunless
- **THEN** the "Sunless" radio SHALL be selected and the other two SHALL NOT be

#### Scenario: All three themes are listed
- **WHEN** the App section is displayed
- **THEN** it SHALL list Parchment, Lumen and Sunless, each with its subtitle

#### Scenario: Selecting a theme card
- **WHEN** the user selects the "Lumen" radio
- **THEN** the "Lumen" radio SHALL become selected and the active theme SHALL change to Lumen
