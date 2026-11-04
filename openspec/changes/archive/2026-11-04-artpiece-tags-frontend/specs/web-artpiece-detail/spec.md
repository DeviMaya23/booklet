# Spec Delta

## ADDED Requirements

### Requirement: View mode — tags
The detail view in view mode SHALL display the artpiece's tags as individual chips beside Artist and Characters, or "—" when the artpiece has no tags.

#### Scenario: Artpiece with tags
- **WHEN** the detail view is open in view mode for an artpiece whose `tags` are `outfit` and `sketch`
- **THEN** a Tags section SHALL show a chip labelled `outfit` and a chip labelled `sketch`

#### Scenario: Artpiece without tags
- **WHEN** the detail view is open in view mode for an artpiece with an empty `tags` array
- **THEN** the Tags section SHALL show "—"

---

### Requirement: Edit mode — tags field
The detail view in edit mode SHALL render a Tags field using the tag picker, full width, below the Artist and Characters row and above Notes. It SHALL be initialised with the artpiece's current tags when edit mode is entered.

#### Scenario: Tags field initialised from the artpiece
- **WHEN** the user enters edit mode on an artpiece with tags `outfit` and `sketch`
- **THEN** the Tags field SHALL show chips for `outfit` and `sketch`

#### Scenario: Add an existing tag
- **WHEN** the user selects an existing tag from the suggestions in edit mode
- **THEN** it SHALL be added as a chip and counted as a pending change

#### Scenario: Create a tag inline in edit mode
- **WHEN** the user types a name matching no existing tag and chooses the create row
- **THEN** the tag SHALL be created, added as a chip, and counted as a pending change; the artpiece itself SHALL NOT be updated until Save

#### Scenario: Remove a tag
- **WHEN** the user removes a chip in edit mode
- **THEN** the tag SHALL be removed from the field and counted as a pending change

#### Scenario: Tags requested in recency order
- **WHEN** the detail view is open in edit mode
- **THEN** the tags request SHALL include `sort=recent` and suggestions SHALL exclude already-selected tags

---

### Requirement: Edit mode — saving tags
Saving in edit mode SHALL send the full selected tag set as `tag_ids` in the `PUT /artpieces/:id` request, replacing the artpiece's tags. After a successful save the tags query SHALL be invalidated.

#### Scenario: Save with changed tags
- **WHEN** the user adds `wip`, removes `sketch` and clicks Save
- **THEN** `PUT /artpieces/:id` is called with `tag_ids` containing the ids of the resulting set, and the view returns to view mode showing the new tags

#### Scenario: Save with all tags removed
- **WHEN** the user removes every tag and clicks Save
- **THEN** `PUT /artpieces/:id` is called with an empty `tag_ids` array

#### Scenario: Save with tags unchanged
- **WHEN** the user changes only the title and clicks Save
- **THEN** `PUT /artpieces/:id` is called with `tag_ids` equal to the artpiece's existing tag ids

#### Scenario: Tags refresh after saving edits
- **WHEN** the user saves artpiece edits
- **THEN** the tags query SHALL be invalidated

#### Scenario: Cancel discards tag edits
- **WHEN** the user changes the tags and clicks Cancel
- **THEN** the changes SHALL be discarded and view mode SHALL show the artpiece's original tags

---

### Requirement: Edit mode — tags in the unsaved-changes check
A difference between the selected tag set and the artpiece's saved tag set, compared by id and ignoring order, SHALL count as an unsaved change for the caption and the route blocker.

#### Scenario: Caption lists tags changed
- **WHEN** the user adds or removes a tag in edit mode
- **THEN** the unsaved changes caption SHALL include "tags changed"

#### Scenario: Reordering or re-adding to the same set is not a change
- **WHEN** the user removes a tag and adds the same tag back
- **THEN** "tags changed" SHALL NOT appear in the caption

#### Scenario: Blocker treats tag-only edits as unsaved
- **WHEN** the user changes only the tags and attempts to navigate away
- **THEN** the route blocker SHALL treat the form as having unsaved changes
