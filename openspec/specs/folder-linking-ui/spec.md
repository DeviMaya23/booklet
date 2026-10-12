# Spec: Folder Linking UI

## Purpose

Defines the UI behaviour for linking Bookleaf public folders to a character from the character create and edit modal, including the folder picker component, degraded-state handling when Bookleaf is unreachable, and diff-based submission semantics.

---

## Requirements

### Requirement: Folder picker in character modal
The character modal SHALL include a Folders section. Each assigned folder SHALL appear as a list row showing a folder icon, the folder name, a share link to the folder in Bookleaf (↗), and a remove (×) button. A "Choose folders" button below the list SHALL open the Choose folders panel for adding new folders. A tooltip (question mark icon) beside the Folders label SHALL explain that the folder list comes from Bookleaf.

#### Scenario: Selected folders render as chips
- **WHEN** the character has existing folder assignments
- **THEN** each folder is displayed as a list row with its name, a Bookleaf share link, and a remove button

#### Scenario: User can remove a folder chip
- **WHEN** the user clicks the remove button on a folder row
- **THEN** that folder is removed from the selection and will not be submitted on save

#### Scenario: User can search and add a folder
- **WHEN** the user opens the Choose folders panel and types in the search bar
- **THEN** the folder list is filtered by name (case-insensitive) and the user can check a folder to add it

#### Scenario: Selected folders are excluded from the add list
- **WHEN** a folder is already assigned to the character
- **THEN** its row in the Choose folders panel shows a checked checkbox

#### Scenario: Tooltip explains Bookleaf source
- **WHEN** the user focuses or hovers the question mark icon beside the Folders label
- **THEN** a tooltip appears explaining that the folder list comes from Bookleaf

#### Scenario: Share link opens Bookleaf folder
- **WHEN** the user clicks the share link (↗) on a folder row
- **THEN** the Bookleaf share URL for that folder (`VITE_BOOKLEAF_HOST/share/:token`) opens in a new tab

#### Scenario: Choose folders button opens the panel
- **WHEN** the user clicks the "Choose folders" button
- **THEN** the Choose folders panel opens alongside the form

---

### Requirement: Folder picker degraded state when Bookleaf is unreachable
If `GET /folders` fails when the Choose folders panel is opened, the panel SHALL display a full-panel error block indicating Bookleaf could not be reached, alongside a Retry action. The error block SHALL note that existing folder assignments are unchanged and that previewing and adding folders require a connection. The main folder list in the form SHALL remain functional with existing folder rows visible and removable. There is no auto-retry — `GET /folders` SHALL be configured with `retry: false`.

#### Scenario: Existing folders render when public-folders fetch fails
- **WHEN** `GET /folders` fails and the character has existing folder assignments
- **THEN** the folder rows in the form remain visible with their remove buttons active

#### Scenario: Add input is disabled when public-folders fetch fails
- **WHEN** the user opens the Choose folders panel and `GET /folders` fails
- **THEN** no folder list is shown and adding or previewing folders is not possible

#### Scenario: Inline error message shown on failure
- **WHEN** the user opens the Choose folders panel and `GET /folders` fails
- **THEN** the panel body displays a full-panel error block ("Couldn't reach Bookleaf") with a Retry action

#### Scenario: Retry triggers a new fetch
- **WHEN** the user clicks Retry inside the error block
- **THEN** a new `GET /folders` request is issued; on success the folder list renders and reconciliation runs

#### Scenario: Unrelated edits are not blocked
- **WHEN** `GET /folders` fails
- **THEN** the user can still edit character name, notes, and avatar, and submit the form

---

### Requirement: Folder name and ID reconciliation on successful fetch
When `GET /folders` succeeds, the modal SHALL silently reconcile stored folders against the live list. No user notice is shown for rename or removal.

- A stored folder whose name differs SHALL have its row name silently updated to the Bookleaf name.
- A stored folder whose ID is absent SHALL be silently removed from the submitted selection. Additionally, a warning row SHALL appear in the Folders section indicating the folder is no longer public in Bookleaf. The warning row SHALL include an explicit Remove button that dismisses it from the UI. The folder is excluded from submitted `folder_ids` regardless of whether the warning row is dismissed.
- A folder added during the current session is not subject to reconciliation.
- A stale cache hit counts as a successful response.

#### Scenario: Stale folder name is silently updated
- **WHEN** `GET /folders` succeeds and a stored folder's name differs from the Bookleaf response
- **THEN** the folder row displays the current Bookleaf name; no indicator is shown

#### Scenario: Folder missing from Bookleaf is silently removed
- **WHEN** `GET /folders` succeeds and a stored folder's ID is absent from the response
- **THEN** a warning row appears in the Folders section indicating the folder is no longer public, and the folder is excluded from the submitted `folder_ids`

#### Scenario: Warning row is dismissed via Remove button
- **WHEN** the user clicks Remove on a warning row
- **THEN** the warning row disappears; the folder remains excluded from submission

#### Scenario: Reconciliation runs after manual retry succeeds
- **WHEN** `GET /folders` initially failed, the user clicked Retry, and the retry succeeded
- **THEN** reconciliation runs against the newly loaded list

#### Scenario: No reconciliation on fetch failure
- **WHEN** `GET /folders` fails (including after retry)
- **THEN** stored folder rows remain as-is with no name or ID changes applied

---

### Requirement: Folder diff submitted on save
The modal SHALL track folder changes as an explicit diff (adds and removes) against the original folder list from the character's GET response. The final `folder_ids` submitted on save SHALL be: original folder IDs plus added IDs minus removed IDs. Folders silently removed during reconciliation are treated as removed in the diff.

#### Scenario: Diff applied correctly on submit
- **WHEN** the original folders are [A, B], the user removes B and adds C via the panel, then saves
- **THEN** the submitted `folder_ids` is [A, C]

#### Scenario: Reconciliation-removed folder excluded from submission
- **WHEN** a stored folder was absent from the Bookleaf response
- **THEN** it is excluded from the submitted `folder_ids` regardless of whether the warning row was dismissed

#### Scenario: Diff submitted even when public-folders fetch failed
- **WHEN** `GET /folders` failed, the user removes an existing folder row, and saves
- **THEN** the submitted `folder_ids` reflects the removal

---

### Requirement: Choose folders panel
The Choose folders panel SHALL list all of the user's public Bookleaf folders. Each row SHALL show a checkbox (checked when assigned), folder name, a Preview button, and a Bookleaf share link. A search bar SHALL filter the list client-side by folder name (case-insensitive). A subtitle SHALL read "Public folders from Bookleaf" alongside a manual refresh button. The panel closes only when the user clicks the X button.

#### Scenario: Panel lists all public folders
- **WHEN** the user opens the Choose folders panel
- **THEN** all public Bookleaf folders are listed, with checkboxes checked for currently assigned folders

#### Scenario: Checking a folder adds it
- **WHEN** the user checks an unchecked folder row
- **THEN** that folder is added to the selection and appears in the form's folder list

#### Scenario: Unchecking a folder removes it
- **WHEN** the user unchecks a checked folder row
- **THEN** that folder is removed from the selection

#### Scenario: Search filters folder list
- **WHEN** the user types in the search bar
- **THEN** only folders whose name contains the search text (case-insensitive) are shown

#### Scenario: Refresh button triggers refetch
- **WHEN** the user clicks the refresh button beside the subtitle
- **THEN** `GET /folders` is re-fetched; on success the list updates and reconciliation runs; on failure the existing list is retained

#### Scenario: Panel closes on X only
- **WHEN** the user clicks the X button in the panel header
- **THEN** the panel closes

---

### Requirement: Folder preview
Each folder row in the Choose folders panel SHALL offer an on-demand image preview fetched from `GET /folders/:id/images`. The FE shows the first 15 images; "Show more" reveals the next 15 from the full response. Only one preview may be open at a time.

#### Scenario: Preview opens on demand
- **WHEN** the user clicks Preview on a folder row
- **THEN** `GET /folders/:id/images` is fetched and images render in a grid below that row

#### Scenario: Skeleton shown while loading
- **WHEN** the preview fetch is in flight
- **THEN** a skeleton grid and a "Loading images…" label are shown

#### Scenario: Initial image count is 15
- **WHEN** the preview loads successfully
- **THEN** up to 15 thumbnails are shown; a "Show more" button appears if more exist

#### Scenario: Show more reveals next 15
- **WHEN** the user clicks "Show more"
- **THEN** the next 15 images from the full response are appended to the grid

#### Scenario: Only one preview open at a time
- **WHEN** the user opens a preview on folder B while folder A's preview is open
- **THEN** folder A's preview collapses and folder B's preview opens

#### Scenario: Hide collapses preview
- **WHEN** the user clicks "Hide" on an open preview
- **THEN** the image grid collapses and the button reverts to "Preview"

---

### Requirement: Choose folders panel states
The panel SHALL handle empty and preview-error states.

- **Empty**: when `GET /folders` returns no folders, the panel SHALL show "No public folders yet" with an "Open Bookleaf" link (navigates to `VITE_BOOKLEAF_HOST`) and a Refresh button.
- **Preview error**: when `GET /folders/:id/images` fails, an inline error and Retry SHALL appear inside that folder's preview area.

#### Scenario: Empty state shown when no public folders
- **WHEN** `GET /folders` returns an empty list
- **THEN** the panel shows "No public folders yet", an "Open Bookleaf" link, and a Refresh button

#### Scenario: Open Bookleaf navigates to Bookleaf root
- **WHEN** the user clicks "Open Bookleaf" in the empty state
- **THEN** `VITE_BOOKLEAF_HOST` opens in a new tab

#### Scenario: Preview error shown inline
- **WHEN** `GET /folders/:id/images` fails for a folder
- **THEN** an inline error message and Retry button appear inside that folder's preview area

#### Scenario: Preview retry re-fetches images
- **WHEN** the user clicks Retry on a preview error
- **THEN** `GET /folders/:id/images` is called again for that folder
