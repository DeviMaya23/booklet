# Spec Delta

## ADDED Requirements

### Requirement: filter_payload includes tag_ids and tag_match_mode
The `filter_payload` object SHALL support `tag_ids` (an array of UUID strings) and `tag_match_mode` (`"all"` or `"any"`). Absent `tag_ids` SHALL read as an empty array and absent `tag_match_mode` as `"all"`. Both are accepted without validation beyond type and are not checked against the user's tags, like `verse_id`.

#### Scenario: Save filter with tag_ids and tag_match_mode
- **WHEN** a user creates a saved filter with filter_payload containing `tag_ids: [<A>, <B>]` and `tag_match_mode: "any"`
- **THEN** the system stores the payload and returns the record with those values unchanged

#### Scenario: Saved filter without tag fields — treated as empty
- **WHEN** a saved filter's filter_payload contains neither `tag_ids` nor `tag_match_mode`
- **THEN** the system treats `tag_ids` as an empty array and `tag_match_mode` as `"all"` when reading that filter

#### Scenario: Deleted tag id remains stored
- **WHEN** a tag listed in a saved filter's `tag_ids` is later deleted
- **THEN** the saved filter record SHALL be unchanged; the stale id is ignored at read time by the client
