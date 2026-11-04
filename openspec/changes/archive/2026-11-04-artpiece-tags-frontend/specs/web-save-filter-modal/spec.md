# Spec Delta

## MODIFIED Requirements

### Requirement: Filter summary chip
The modal SHALL display a read-only chip summarising the active filter state in human-readable form, including the tags and the verse when active (e.g. "Artist: Chev · Characters: Tanelle, Jielle (any) · Tags: outfit, sketch (all) · Verse: Elowen").

#### Scenario: Summary reflects active filters
- **WHEN** the modal opens with artist "Chev" and characters "Tanelle" and "Jielle" with match mode "any"
- **THEN** the chip SHALL display text representing those active filters

#### Scenario: Summary includes active verse
- **WHEN** the modal opens with verse "Elowen" active (alone or alongside other filters)
- **THEN** the chip SHALL include the verse's name in the summary text

#### Scenario: Summary includes active tags
- **WHEN** the modal opens with tags "outfit" and "sketch" active with tag match mode "all" (alone or alongside other filters)
- **THEN** the chip SHALL include the tag names and the tag match mode in the summary text

---

### Requirement: Save action
Clicking Save SHALL POST to `/saved_filters` with the name, `thumbnail_artpiece_id`, and `filter_payload`, where `filter_payload` includes `tag_ids` and `tag_match_mode`. On success the modal closes and the caller receives the created saved filter. On a non-422 error a toast fires and the modal stays open.

#### Scenario: Successful save closes modal
- **WHEN** the server returns 201
- **THEN** the modal SHALL close and a success toast SHALL appear

#### Scenario: Payload carries tags
- **WHEN** the user saves with tags `outfit` and `sketch` selected and tag match mode "any"
- **THEN** the POST body's `filter_payload` SHALL contain `tag_ids` with both ids and `tag_match_mode: "any"`

#### Scenario: Payload without tags
- **WHEN** the user saves with no tags selected
- **THEN** the POST body's `filter_payload` SHALL contain `tag_ids: []`

#### Scenario: Only resolved tags are saved
- **WHEN** the URL carried a tag id that no longer exists and the user saves the filter
- **THEN** `tag_ids` SHALL contain only the ids of tags that currently exist

#### Scenario: Server error shows toast
- **WHEN** the server returns a non-422 error
- **THEN** the modal SHALL remain open and an error toast SHALL appear

#### Scenario: Cancel discards without saving
- **WHEN** the user clicks Cancel
- **THEN** the modal SHALL close without POSTing to the server
