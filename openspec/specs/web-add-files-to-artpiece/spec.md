## Purpose

Defines the "Add to Existing Artpiece" modal — the flow that lets a user attach a set of already-selected inbox files to an existing artpiece by searching for it, previewing it, and confirming.

## Requirements

### Requirement: Add to existing artpiece modal
The system SHALL provide a modal that lets the user attach a set of already-selected inbox files to an existing artpiece by searching for it, previewing it, and confirming.

#### Scenario: Modal opens with selected files
- **WHEN** the user opens the "Add to Existing Artpiece" modal with N files selected
- **THEN** the modal header reads "Add N file(s) to artpiece" and the top section shows a thumbnail strip of the selected files

#### Scenario: Thumbnail strip — file with thumbnail
- **WHEN** a selected file has a `thumbnail_url`
- **THEN** the tile shows the thumbnail image with the file's name truncated below it (full name on hover)

#### Scenario: Thumbnail strip — file without thumbnail
- **WHEN** a selected file has no `thumbnail_url`
- **THEN** the tile shows a fallback file icon with the file's name truncated below it (full name on hover)

#### Scenario: Thumbnail strip scroll
- **WHEN** the selected files produce more than two rows of thumbnails
- **THEN** the strip container becomes vertically scrollable at a fixed max-height; no horizontal scroll is introduced

#### Scenario: Search by title — results appear
- **WHEN** the user types in the "search by title" input
- **THEN** the artpiece list is filtered client-side (case-insensitive substring match on `title`) and results are shown below the input; no network request is made

#### Scenario: Search by title — no results
- **WHEN** the search term matches no artpiece titles
- **THEN** the results area shows an empty state and no artpiece preview is shown

#### Scenario: Filters section collapsed by default
- **WHEN** the modal opens
- **THEN** the Filters section is collapsed and shows only the "Filters ▽" toggle

#### Scenario: Filters section expanded
- **WHEN** the user clicks the "Filters ▽" toggle
- **THEN** the filters section expands to show an Artist dropdown and a Characters chip multi-select

#### Scenario: Artist filter
- **WHEN** the user selects an artist in the Artist filter
- **THEN** the artpiece list is filtered client-side to include only artpieces by that artist

#### Scenario: Characters filter
- **WHEN** the user selects one or more characters in the Characters filter
- **THEN** the artpiece list is filtered client-side to include only artpieces that include all selected characters

#### Scenario: Combined filters
- **WHEN** the user has both a title search and one or more filters active simultaneously
- **THEN** all filters apply together (intersection) client-side

#### Scenario: Picking a search result
- **WHEN** the user clicks a result in the artpiece list
- **THEN** the system fetches `GET /artpieces/:id`, shows the picked-artpiece preview (cover thumbnail if any, title, artist, character(s)), and holds the current file list for merge

#### Scenario: Picked artpiece preview — no cover
- **WHEN** the picked artpiece has no cover thumbnail
- **THEN** the cover area in the preview shows a fallback placeholder

#### Scenario: Changing the pick
- **WHEN** the user has already picked an artpiece and then clicks a different search result
- **THEN** the preview updates to the newly picked artpiece

#### Scenario: Save — success
- **WHEN** the user clicks Save with a picked artpiece
- **THEN** the system calls `PUT /artpieces/:id/files` with the artpiece's existing file IDs merged with the selected inbox file IDs, shows a success toast, clears the inbox selection, and closes the modal

#### Scenario: Save — file already attached to a different artpiece
- **WHEN** the backend returns HTTP 422 on save (e.g., a selected file was attached elsewhere since the modal opened)
- **THEN** the modal shows an error toast and remains open

#### Scenario: Save disabled without pick
- **WHEN** no artpiece has been picked
- **THEN** the Save button is disabled
