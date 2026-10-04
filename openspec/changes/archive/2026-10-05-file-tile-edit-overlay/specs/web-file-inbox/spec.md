## ADDED Requirements

### Requirement: File name label beneath tile
The system SHALL display the file's name as a read-only, truncated single-line label beneath each tile in the inbox grid.

#### Scenario: File with a name
- **WHEN** a file in the grid has a non-null `name`
- **THEN** the name is displayed in a single truncated line below the tile image

#### Scenario: File with no name
- **WHEN** a file in the grid has a null `name`
- **THEN** no label is shown beneath the tile

---

### Requirement: Double-click tile opens edit overlay
The system SHALL open the file edit overlay when the user double-clicks a tile. Double-click does not affect the tile's selection state.

#### Scenario: Double-clicking a tile
- **WHEN** a user double-clicks a file tile
- **THEN** the edit overlay opens for that file; the tile's selection state is unchanged

## MODIFIED Requirements

### Requirement: Per-tile actions menu
The system SHALL show a trash icon on each tile on hover for the Delete action. The "..." dropdown menu and "View Detail" item are removed.

#### Scenario: Hovering a tile
- **WHEN** the user hovers over a file tile
- **THEN** a trash icon appears on the tile

#### Scenario: Clicking the trash icon
- **WHEN** the user clicks the trash icon on a tile
- **THEN** a confirmation dialog appears for that single file; on confirm, the file is deleted
