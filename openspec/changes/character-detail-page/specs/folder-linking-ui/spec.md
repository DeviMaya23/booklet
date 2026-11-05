# Spec Delta

## ADDED Requirements

### Requirement: Folder rules apply on the character detail page
In edit mode the character detail page SHALL apply the same folder behavior as the create modal: degraded-state handling when `GET /folders` fails, silent reconciliation of renamed and no-longer-public folders against the live list, and a submitted `folder_ids` equal to the original IDs plus added IDs minus removed IDs.

#### Scenario: Reconciliation runs on the detail page
- **WHEN** the detail page is in edit mode, `GET /folders` succeeds and an attached folder's ID is absent
- **THEN** a "No longer public in Bookleaf." warning row appears and the folder is excluded from the submitted `folder_ids`

#### Scenario: Renamed folder shows the live name
- **WHEN** an attached folder's stored name differs from the Bookleaf name
- **THEN** its chip shows the Bookleaf name

#### Scenario: Diff applied on save from the detail page
- **WHEN** the original folders are [A, B], the user removes B and adds C in the drawer, and saves
- **THEN** the submitted `folder_ids` is [A, C]

#### Scenario: Degraded state in the drawer
- **WHEN** `GET /folders` fails and the user opens the drawer
- **THEN** the panel shows the "Couldn't reach Bookleaf" block with Retry, and existing chips stay removable

### Requirement: Folder preview paging and links
Each folder row in the Choose folders panel SHALL offer an on-demand image preview fetched from `GET /folders/:id/images`. The FE shows the first 24 images; "Show more" reveals the next 24 from the full response. A caption ("Showing the first N" while more remain, otherwise "Showing all N") and a "View all in Bookleaf" link to the folder's share URL SHALL accompany the grid. Only one preview may be open at a time.

#### Scenario: Preview opens on demand
- **WHEN** the user clicks Preview on a folder row
- **THEN** `GET /folders/:id/images` is fetched and images render in a grid below that row

#### Scenario: Skeleton shown while loading
- **WHEN** the preview fetch is in flight
- **THEN** a skeleton grid and a "Loading images…" label are shown

#### Scenario: Initial image count is 24
- **WHEN** the preview loads successfully
- **THEN** up to 24 thumbnails are shown with a "Showing the first 24" caption; a "Show more" button appears if more exist

#### Scenario: Caption reads "Showing all" when nothing remains
- **WHEN** every image in the folder is visible
- **THEN** the caption reads "Showing all N" and no "Show more" button is shown

#### Scenario: Show more reveals next 24
- **WHEN** the user clicks "Show more"
- **THEN** the next 24 images from the full response are appended to the grid

#### Scenario: View all opens the folder in Bookleaf
- **WHEN** the user clicks "View all in Bookleaf"
- **THEN** the folder's Bookleaf share URL opens in a new tab

#### Scenario: Only one preview open at a time
- **WHEN** the user opens a preview on folder B while folder A's preview is open
- **THEN** folder A's preview collapses and folder B's preview opens

#### Scenario: Hide collapses preview
- **WHEN** the user clicks "Hide" on an open preview
- **THEN** the image grid collapses and the button reverts to "Preview"

## MODIFIED Requirements

### Requirement: Choose folders panel
The Choose folders panel SHALL list all of the user's public Bookleaf folders. Each row SHALL show a checkbox (checked when assigned), folder name, a Preview button, and a Bookleaf share link. A search bar SHALL filter the list client-side by folder name (case-insensitive). A subtitle SHALL read "Public folders from Bookleaf" alongside an "Updated <relative time>" label for the loaded list and a manual refresh button. The panel closes only when the user clicks the X button (or, in the detail-page drawer, presses Escape).

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

#### Scenario: Updated label reflects the last successful fetch
- **WHEN** the folder list was last fetched 3 minutes ago
- **THEN** the panel shows an "Updated" label with a relative time such as "3 minutes ago"

#### Scenario: Refresh button triggers refetch
- **WHEN** the user clicks the refresh button beside the subtitle
- **THEN** `GET /folders` is re-fetched; on success the list updates, the Updated label resets and reconciliation runs; on failure the existing list is retained

#### Scenario: Panel closes on X only
- **WHEN** the user clicks the X button in the panel header
- **THEN** the panel closes

## REMOVED Requirements

### Requirement: Folder preview
**Reason**: Page size changed from 15 to 24 and a caption and "View all in Bookleaf" link were added; replaced by "Folder preview paging and links".
**Migration**: See "Folder preview paging and links" in this capability. Preview-on-demand, one-open-at-a-time, skeleton and Hide behavior are carried over unchanged.
