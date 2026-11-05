# Spec Delta

## Purpose

Defines the Tags page of the web app, where a user reviews all of their tags and renames, merges or deletes them. It lets a single-user library be cleaned up without touching individual artpieces.

## ADDED Requirements

### Requirement: Tags page route and navigation
The app SHALL provide a Tags page at `/app/tags`, reachable from the "Tags" item in the sidebar's MANAGE group.

#### Scenario: Open from the sidebar
- **WHEN** an authenticated user clicks "Tags" in the sidebar
- **THEN** the app navigates to `/app/tags` and renders the Tags page with the Tags nav row active

#### Scenario: Direct navigation
- **WHEN** an authenticated user opens `/app/tags` directly
- **THEN** the Tags page renders

### Requirement: Tag list with usage
The Tags page SHALL list every tag the user owns, one row per tag. Each row SHALL show the tag name, a usage label "N artpieces" (using the singular "1 artpiece" when N is 1), an edit action and a delete action. The data SHALL be loaded with usage counts requested.

#### Scenario: Rows render
- **WHEN** the user has tags `outfit` (22 artpieces) and `sketch` (14 artpieces)
- **THEN** two rows render showing the names with "22 artpieces" and "14 artpieces"

#### Scenario: Unused tag is badged
- **WHEN** a tag is used on 0 artpieces
- **THEN** its row shows an "Unused" badge next to the name and the usage label "0 artpieces" without a link

#### Scenario: Singular label
- **WHEN** a tag is used on exactly 1 artpiece
- **THEN** its usage label reads "1 artpiece"

### Requirement: Usage label links to filtered gallery
For a tag used on at least one artpiece, the usage label SHALL be a link to the artpieces gallery filtered by that tag.

#### Scenario: Follow the usage link
- **WHEN** the user clicks "14 artpieces" on the row for tag `sketch`
- **THEN** the app navigates to `/app/artpieces?tags=<sketch tag id>`

### Requirement: Tag search
The Tags page SHALL provide a search box that filters the displayed tags to those whose name contains the typed text, ignoring case. Filtering SHALL happen on the already-loaded list.

#### Scenario: Filter by substring
- **WHEN** the user types `sket` with tags `sketch`, `sketchbook` and `chibi` present
- **THEN** only `sketch` and `sketchbook` are shown

#### Scenario: No matches
- **WHEN** the search text matches no tag but the user has tags
- **THEN** no rows are shown and no "No tags yet" empty state is shown

### Requirement: Tag sort
The Tags page SHALL provide a sort dropdown with the options "Recently used" and "Name (A–Z)". "Recently used" SHALL be the default and SHALL use the order returned by the API for `sort=recent`. "Name (A–Z)" SHALL use the API's default order.

#### Scenario: Default sort is recently used
- **WHEN** the user opens the Tags page
- **THEN** the dropdown shows "Recently used" and rows appear in the `sort=recent` order

#### Scenario: Switch to name order
- **WHEN** the user selects "Name (A–Z)"
- **THEN** rows appear in ascending name order

#### Scenario: Search and sort combine
- **WHEN** the user has typed search text and changes the sort
- **THEN** the filtered rows are shown in the newly selected order

### Requirement: Empty state
When the user has no tags, the Tags page SHALL show an empty state reading "No tags yet" with the hint that tags added from an artpiece's edit page appear there, in place of the list.

#### Scenario: No tags
- **WHEN** a user with no tags opens the Tags page
- **THEN** the "No tags yet" empty state is shown and no rows render

### Requirement: Inline rename
Activating a row's edit action SHALL replace that row's content with a text input pre-filled with the current name and focused, a save control and a cancel control, and the hint "Enter to save · Esc to cancel". Enter or the save control SHALL submit the rename; Esc or the cancel control SHALL discard it. Only one row SHALL be in edit mode at a time.

#### Scenario: Enter edit mode
- **WHEN** the user clicks the edit action on the row for `outfit`
- **THEN** the row shows a focused input containing `outfit`, the save and cancel controls and the hint text

#### Scenario: Save a rename
- **WHEN** the user changes the input to `outfits` and presses Enter and no other tag is named `outfits`
- **THEN** the tag is renamed through the API, the row leaves edit mode and shows `outfits`

#### Scenario: Cancel with Escape
- **WHEN** the user presses Esc while editing
- **THEN** the row leaves edit mode with the original name and no request is sent

#### Scenario: Unchanged or blank name
- **WHEN** the user submits the original name unchanged, or a blank name
- **THEN** no rename request is sent; an unchanged name leaves edit mode and a blank name keeps the input open

#### Scenario: Rename failure
- **WHEN** the rename request fails
- **THEN** the row stays in edit mode and an error is shown

### Requirement: Merge on rename collision
When the submitted name, compared case-insensitively after trimming, equals the name of another of the user's tags, the app SHALL NOT call the rename endpoint and SHALL show a "Merge into “<target>”?" confirmation. Confirming SHALL merge through the merge endpoint; cancelling SHALL leave the row in edit mode.

#### Scenario: Dialog explains the consequences
- **WHEN** the merge dialog is shown for renaming `sketches` into `sketch`
- **THEN** it states that artpieces tagged `sketches` get `sketch` instead, that artpieces having both keep just `sketch`, that `sketches` is deleted, and that this can't be undone

#### Scenario: Collision opens merge dialog
- **WHEN** the user renames `sketches` to `sketch` while a tag `sketch` exists
- **THEN** the "Merge into “sketch”?" dialog opens and no request has been sent

#### Scenario: Case-insensitive collision
- **WHEN** the user renames `sketches` to `Sketch` while a tag `sketch` exists
- **THEN** the merge dialog for `sketch` opens

#### Scenario: Confirm merge
- **WHEN** the user clicks "Merge tags"
- **THEN** the app calls the merge endpoint with `sketches` as source and `sketch` as target, the dialog closes, and the list refreshes without `sketches`

#### Scenario: Cancel merge
- **WHEN** the user clicks "Cancel" in the merge dialog
- **THEN** the dialog closes, no request is sent, and the row is still in edit mode

#### Scenario: Case-only rename is not a collision
- **WHEN** the user renames `wip` to `WIP` and no other tag is named `wip`
- **THEN** the normal rename is submitted, not the merge flow

### Requirement: Delete tag with confirmation
Activating a row's delete action SHALL open a "Delete “<name>”?" confirmation stating how many artpieces use the tag, that they will lose it, and that nothing else about them changes. Confirming SHALL delete the tag through the API; cancelling SHALL close the dialog without a request.

#### Scenario: Confirm deletion of a used tag
- **WHEN** the user clicks delete on `sketch` (14 artpieces) and then "Delete"
- **THEN** the dialog text states 14 artpieces use the tag, the tag is deleted through the API, and its row disappears

#### Scenario: Delete an unused tag
- **WHEN** the user clicks delete on a tag used on 0 artpieces
- **THEN** the dialog states that no artpieces use the tag

#### Scenario: Cancel deletion
- **WHEN** the user clicks "Cancel" in the delete dialog
- **THEN** the dialog closes and the tag remains

### Requirement: Tag changes refresh dependent tag data
After a successful rename, merge or delete, the app SHALL refresh all cached tag lists, so that tag pickers and filters elsewhere in the app reflect the change.

#### Scenario: Picker reflects a merge
- **WHEN** the user merges `sketches` into `sketch` and then opens an artpiece's tag picker
- **THEN** `sketches` is not offered
