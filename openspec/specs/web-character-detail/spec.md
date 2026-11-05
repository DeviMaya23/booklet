# web-character-detail Specification

## Purpose

Defines the character detail page at `/app/characters/:id`: a place to browse a character's identity, attached Bookleaf folders and the artpieces and references it appears in, and to edit or delete the character in place.

## Requirements

### Requirement: Character detail is accessible via URL
The system SHALL provide a route `/app/characters/:id` that renders the detail page for the given character. Clicking a character card on `/app/characters` SHALL navigate to this route. The Characters sidebar item SHALL remain active on this route.

#### Scenario: Detail page loads the character
- **WHEN** a user navigates to `/app/characters/:id` with a valid character ID
- **THEN** the page fetches `GET /characters/:id` and renders that character in view mode

#### Scenario: Sidebar highlights Characters on the detail route
- **WHEN** a user is on `/app/characters/:id`
- **THEN** the Characters sidebar item appears active

#### Scenario: Unknown character shows an error state
- **WHEN** `GET /characters/:id` fails or the character does not exist
- **THEN** the page shows "Failed to load character." with a back link to `/app/characters`

#### Scenario: Loading state
- **WHEN** the character request is in flight
- **THEN** the page shows a loading indicator and the back link

### Requirement: View mode header
In view mode the page SHALL show a "← Characters" back link to `/app/characters`, the character's avatar tile, the name, an Edit button, and the character's notes. When the character has a verse, a chip with the verse name SHALL appear beside the name. When `avatar_url` is null the tile SHALL show a placeholder.

#### Scenario: Header shows character details
- **WHEN** the page renders a character with avatar, verse "Coffee Shop AU" and notes "Apothecary from Dyffrynne"
- **THEN** the avatar, the name, a "Coffee Shop AU" chip and the notes text are shown

#### Scenario: No verse hides the chip
- **WHEN** the character has no verse
- **THEN** no verse chip is shown

#### Scenario: No notes hides the notes line
- **WHEN** the character has null notes
- **THEN** no notes text is shown

#### Scenario: No avatar shows a placeholder
- **WHEN** the character's `avatar_url` is null
- **THEN** the avatar tile shows a placeholder image icon

#### Scenario: Back link returns to the list
- **WHEN** the user clicks "← Characters" in view mode
- **THEN** the browser navigates to `/app/characters`

### Requirement: View mode folders row
The header SHALL include a Folders row. Each attached folder SHALL render as a chip with its name. When the live Bookleaf folder list is available and contains the folder, the chip SHALL link to `VITE_BOOKLEAF_HOST/share/:token` in a new tab. When the character has no folders the row SHALL read "No folders attached".

#### Scenario: Attached folders render as chips
- **WHEN** the character has folders "Plant" and "Cottagecore"
- **THEN** the Folders row shows a chip for each name

#### Scenario: Chips link to Bookleaf when the live list is available
- **WHEN** `GET /folders` succeeds and includes an attached folder
- **THEN** that folder's chip links to its Bookleaf share URL in a new tab

#### Scenario: Chips degrade to plain names when Bookleaf is unreachable
- **WHEN** `GET /folders` fails
- **THEN** the folder chips still render with their names, without links, and no error is shown

#### Scenario: No folders attached
- **WHEN** the character has no folders
- **THEN** the Folders row shows "No folders attached"

### Requirement: Artpieces and References tabs
Below the header the page SHALL show an Artpieces tab and a References tab, each with a count. Both SHALL list artpieces whose `characters` include this character, taken from `GET /artpieces`. Artpieces with `is_reference: false` SHALL appear only under Artpieces; artpieces with `is_reference: true` SHALL appear only under References. Artpieces is selected by default.

#### Scenario: Artpieces are split by reference flag
- **WHEN** the character appears in two non-reference artpieces and one reference artpiece
- **THEN** the Artpieces tab shows count 2 and the two artpieces, and the References tab shows count 1 and the reference artpiece

#### Scenario: Other characters' artpieces are excluded
- **WHEN** an artpiece does not list this character
- **THEN** it appears in neither tab

#### Scenario: Reference artpieces are not shown under Artpieces
- **WHEN** the Artpieces tab is selected
- **THEN** no artpiece with `is_reference: true` is shown

#### Scenario: Loading is not shown as empty
- **WHEN** `GET /artpieces` is still in flight
- **THEN** the tab panel shows "Loading…", the tab counts are hidden, and no empty-state block is shown

#### Scenario: Failure is not shown as empty
- **WHEN** `GET /artpieces` fails and there is no cached data
- **THEN** the tab panel shows "Failed to load artpieces.", the tab counts are hidden, and no empty-state block is shown

#### Scenario: Switching tabs
- **WHEN** the user selects the References tab
- **THEN** the references grid is shown in place of the artpieces grid

### Requirement: Tab empty states
When the selected tab has no items, the page SHALL show a dashed empty-state block that names the character and tells the user to add the character in an artpiece's Characters field. The copy SHALL NOT assume the character's pronouns.

#### Scenario: Empty Artpieces tab
- **WHEN** the character has no non-reference artpieces and the Artpieces tab is selected
- **THEN** the empty-state block shows "No artpieces with <name> yet"

#### Scenario: Empty References tab
- **WHEN** the character has no reference artpieces and the References tab is selected
- **THEN** the empty-state block shows "No references with <name> yet"

### Requirement: Opening artpieces from the tabs
Clicking a tile in either tab SHALL open the file viewer on that tab's artpieces, as the gallery does. Using a tile's info control SHALL navigate to `/app/artpieces/:id`.

#### Scenario: Tile click opens the viewer
- **WHEN** the user clicks an artpiece tile in the Artpieces tab
- **THEN** the file viewer opens on that artpiece among the tab's artpieces

#### Scenario: Info control opens the artpiece detail
- **WHEN** the user activates a tile's info control
- **THEN** the browser navigates to `/app/artpieces/:id` for that artpiece

### Requirement: Edit mode header
Clicking Edit SHALL switch the page to edit mode in place. The header SHALL show a back link labelled with the character's name ("← Tanelle") that returns to view mode, a name input, and Cancel and Save buttons. Save SHALL be disabled while the name is blank or a save is in progress.

#### Scenario: Edit switches the page to edit mode
- **WHEN** the user clicks Edit
- **THEN** the name input, Cancel, Save, verse, notes, avatar and folder controls appear, pre-filled from the character

#### Scenario: Blank name blocks Save
- **WHEN** the name input is empty or whitespace
- **THEN** Save is disabled and no request is made

#### Scenario: Cancel discards edits
- **WHEN** the user clicks Cancel
- **THEN** edit state is discarded and the page returns to view mode showing the original values

### Requirement: Unsaved changes summary
In edit mode, while any field differs from the saved character, the page SHALL show an amber-dotted line "Unsaved changes: …" listing what changed, such as "name changed", "verse changed", "notes changed", "avatar changed", "avatar removed", "N folder(s) added" and "N folder(s) removed". Folders dropped automatically because Bookleaf no longer lists them SHALL NOT count as changes. The line SHALL be hidden when nothing differs.

#### Scenario: Summary lists changes
- **WHEN** the user changes the verse and removes one folder
- **THEN** the line reads "Unsaved changes: verse changed, 1 folder removed"

#### Scenario: Reverting a change clears it
- **WHEN** the user changes the verse and then selects the original verse again
- **THEN** "verse changed" no longer appears and the line is hidden if nothing else differs

#### Scenario: Reconciliation alone is not a change
- **WHEN** a stored folder is no longer public in Bookleaf and the user has changed nothing
- **THEN** no "Unsaved changes" line is shown and leaving the page does not prompt

#### Scenario: User changes are counted next to a dropped folder
- **WHEN** a stored folder is no longer public in Bookleaf and the user removes another folder
- **THEN** the summary reads "1 folder removed", counting only the user's removal

### Requirement: Leaving with unsaved changes
In edit mode with unsaved changes, using the back link or navigating away (browser navigation or sidebar) SHALL open a "Discard changes?" dialog with "Keep editing" and "Discard". Discard SHALL drop the edits and complete the action (return to view mode, or proceed with the navigation); Keep editing SHALL keep the edits. Without unsaved changes there SHALL be no dialog.

#### Scenario: Discard dialog on back link
- **WHEN** the user has unsaved changes and clicks the back link
- **THEN** the "Discard changes?" dialog opens and the page stays in edit mode

#### Scenario: Discard from the back link
- **WHEN** the user confirms Discard after clicking the back link
- **THEN** the edits are dropped and the page returns to view mode

#### Scenario: Keep editing stays
- **WHEN** the user chooses Keep editing
- **THEN** the dialog closes and the edits remain

#### Scenario: Sidebar or browser navigation is guarded
- **WHEN** the user has unsaved changes and navigates to another page in the app
- **THEN** the "Discard changes?" dialog opens, and confirming Discard completes the navigation

#### Scenario: No dialog without changes
- **WHEN** the user is in edit mode with no changes and clicks the back link or navigates away
- **THEN** the action completes immediately with no dialog

### Requirement: Edit mode fields
Edit mode SHALL provide a Verse combobox with a clear control and a Notes textarea. The Verse field SHALL offer the user's existing verses and SHALL allow clearing to no verse.

#### Scenario: Change verse
- **WHEN** the user picks a different verse and saves
- **THEN** the request carries the new `verse_id`

#### Scenario: Clear verse
- **WHEN** the user clears the Verse field and saves
- **THEN** the request carries `verse_id: null`

#### Scenario: Clear notes
- **WHEN** the user empties the Notes textarea and saves
- **THEN** the request carries `notes: null`

### Requirement: Avatar editing on the detail page
In edit mode the avatar tile SHALL show a "Change" overlay when an avatar exists and a remove (X) control. Clicking the tile, dropping an image on it, or pasting an image with no text field focused SHALL open the crop dialog. A confirmed crop SHALL preview locally; removing SHALL revert the tile to the placeholder. Neither SHALL be sent before Save.

#### Scenario: Change opens the picker and crop dialog
- **WHEN** the user clicks the avatar tile and chooses an image within the size limit
- **THEN** the crop dialog opens, and on confirm the tile shows the cropped preview

#### Scenario: Drop and paste are accepted
- **WHEN** the user drops an image on the tile, or pastes an image with no text field focused
- **THEN** the crop dialog opens with that image

#### Scenario: Remove avatar locally
- **WHEN** the user clicks the X on the avatar tile
- **THEN** the tile shows the placeholder and no request is made until Save

#### Scenario: Oversized image rejected
- **WHEN** the chosen image exceeds the image size limit
- **THEN** a toast names the limit and the crop dialog does not open

### Requirement: Saving a character from the detail page
On Save the page SHALL: upload a newly cropped avatar (init → PUT to the presigned URL → complete) or `DELETE /characters/:id/avatar` if it was removed; then call `PUT /characters/:id` with name, notes, `folder_ids` and `verse_id`. On success it SHALL show a success toast and return to view mode with the saved values. On failure it SHALL show an error toast and stay in edit mode with edits intact.

#### Scenario: Metadata-only save
- **WHEN** the user edits name or notes and saves
- **THEN** `PUT /characters/:id` is called, a success toast is shown, and view mode shows the new values

#### Scenario: Saved values show immediately
- **WHEN** the update succeeds and the page returns to view mode
- **THEN** view mode shows the saved values straight away, never the previous ones

#### Scenario: Save with new avatar
- **WHEN** the user crops a new avatar and saves
- **THEN** the avatar upload sequence completes before `PUT /characters/:id`

#### Scenario: Save after removing the avatar
- **WHEN** the user removes the avatar and saves
- **THEN** `DELETE /characters/:id/avatar` is called before `PUT /characters/:id`

#### Scenario: Avatar upload failure
- **WHEN** the avatar upload fails
- **THEN** an error toast is shown, no `PUT` is sent, and the page stays in edit mode

#### Scenario: Update failure
- **WHEN** `PUT /characters/:id` fails
- **THEN** an error toast is shown and the page stays in edit mode with edits intact

### Requirement: Folders in edit mode
Edit mode SHALL show a Folders section with a count, a description that the folders are Bookleaf folders whose images belong to the character, one chip per selected folder (name, open-in-Bookleaf link when a token is known, remove control), a warning row for each folder no longer public in Bookleaf, and a "Choose folders" button.

#### Scenario: Chips show remove and Bookleaf link
- **WHEN** the character has folders and the live list is available
- **THEN** each chip shows its name, an open-in-Bookleaf link and a remove control

#### Scenario: Removing a chip
- **WHEN** the user clicks remove on a folder chip
- **THEN** the folder leaves the selection and is excluded from `folder_ids` on Save

#### Scenario: No longer public folder
- **WHEN** an attached folder is absent from the live Bookleaf list
- **THEN** a warning row names it with "No longer public in Bookleaf." and a Remove button, and it is excluded from `folder_ids` on Save

#### Scenario: Bookleaf unreachable
- **WHEN** `GET /folders` fails
- **THEN** existing chips remain visible and removable, without links

### Requirement: Choose folders drawer
"Choose folders" SHALL open the folder picker in a right-side drawer over the page. The drawer SHALL contain the Choose folders panel, SHALL close with its X button or Escape, and SHALL leave the page's edit state intact when closed.

#### Scenario: Open and close the drawer
- **WHEN** the user clicks "Choose folders" and then the drawer's X button
- **THEN** the drawer opens over the page and then closes, and the page remains in edit mode

#### Scenario: Selections made in the drawer show in the page
- **WHEN** the user checks a folder in the drawer and closes it
- **THEN** a chip for that folder appears in the Folders section and "1 folder added" appears in the unsaved-changes line

### Requirement: Delete a character from the detail page
Edit mode SHALL include a "Delete character" action. It SHALL open a confirmation dialog; on confirm it SHALL call `DELETE /characters/:id`, show a success toast and navigate to `/app/characters`. On cancel nothing SHALL change.

#### Scenario: Delete confirmed
- **WHEN** the user confirms deletion
- **THEN** `DELETE /characters/:id` is called, a success toast is shown, and the browser navigates to `/app/characters`

#### Scenario: Delete cancelled
- **WHEN** the user dismisses the confirmation dialog
- **THEN** the dialog closes and the character is unchanged

#### Scenario: Delete failure
- **WHEN** `DELETE /characters/:id` fails
- **THEN** an error toast is shown and the page stays on the character
