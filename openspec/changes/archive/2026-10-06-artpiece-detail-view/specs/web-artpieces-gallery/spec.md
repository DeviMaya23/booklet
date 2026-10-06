## ADDED Requirements

### Requirement: Open artpiece detail via double-click
The gallery card SHALL support double-click to open the artpiece detail view.

#### Scenario: Double-click opens detail view
- **WHEN** the user double-clicks an artpiece card in the gallery
- **THEN** the artpiece detail view SHALL open for that artpiece

---

## MODIFIED Requirements

### Requirement: Delete artpiece from gallery
The gallery SHALL allow the user to delete an artpiece via the card's overflow menu. A confirmation dialog with an "Also delete attached files" checkbox SHALL appear. When the checkbox is checked, attached files are permanently deleted along with the artpiece; when unchecked, files remain in the inbox with no artpiece association.

#### Scenario: Delete triggers confirmation dialog with checkbox
- **WHEN** the user clicks "Delete" in a card's overflow menu
- **THEN** the `DeleteArtpieceDialog` SHALL open with the "Also delete attached files" checkbox (generic label, no file count) and a Cancel and Delete button

#### Scenario: Confirmed deletion without delete_files
- **WHEN** the user confirms deletion with the checkbox unchecked
- **THEN** `DELETE /artpieces/:id` SHALL be called; on success the card SHALL be removed from the gallery and a toast SHALL show "Artpiece deleted"

#### Scenario: Confirmed deletion with delete_files
- **WHEN** the user confirms deletion with the checkbox checked
- **THEN** `DELETE /artpieces/:id?delete_files=true` SHALL be called; on success the card SHALL be removed from the gallery and a toast SHALL show "Artpiece deleted"

#### Scenario: Cancelled deletion does nothing
- **WHEN** the user dismisses the confirmation dialog
- **THEN** no API call is made and the gallery is unchanged

#### Scenario: Deletion failure shows error toast
- **WHEN** the delete call fails
- **THEN** the card SHALL remain in the gallery and a toast SHALL show "Failed to delete artpiece"
