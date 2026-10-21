# web-thumbnail-tile Specification

## Purpose

Defines the shared `ThumbnailTile` component — a navigation-first square thumbnail card with a label, optional sublabel, and an optional contextual action menu. Used wherever a clickable thumbnail navigates to another page.

## Requirements

### Requirement: Thumbnail image or placeholder
The component SHALL render a square image from the provided URL. When no URL is provided or the image fails to load, the component SHALL render a muted placeholder in its place.

#### Scenario: Image renders when URL provided
- **WHEN** a thumbnail URL is provided
- **THEN** the tile SHALL display the image filling the square area

#### Scenario: Placeholder renders when no URL
- **WHEN** no thumbnail URL is provided
- **THEN** the tile SHALL render a muted placeholder icon in the square area

---

### Requirement: Navigation primary action
The tile's square image area SHALL act as a navigable link. Clicking it SHALL navigate to the provided `href`.

#### Scenario: Click navigates
- **WHEN** a user clicks the tile's image area
- **THEN** the app SHALL navigate to the tile's configured `href`

---

### Requirement: Label and sublabel
The tile SHALL display a label below the image. An optional sublabel SHALL render below the label in a muted style when provided.

#### Scenario: Label renders
- **WHEN** a tile is rendered with a label
- **THEN** the label SHALL appear below the image area

#### Scenario: Sublabel renders when provided
- **WHEN** a sublabel is provided
- **THEN** the sublabel SHALL render below the label in a muted style

#### Scenario: Sublabel absent when not provided
- **WHEN** no sublabel is provided
- **THEN** no sublabel element SHALL be rendered

---

### Requirement: Optional action menu
When `menuItems` are provided, the tile SHALL render a contextual action button (visible on hover) that opens a dropdown menu with those items. When `menuItems` is absent or empty, no action button SHALL be rendered.

#### Scenario: Action button absent when no menu items
- **WHEN** the tile is rendered with no `menuItems`
- **THEN** no action button SHALL be present on the tile

#### Scenario: Action button opens menu
- **WHEN** the tile has `menuItems` and the user activates the action button
- **THEN** a dropdown menu SHALL open listing those items

#### Scenario: Destructive item renders in destructive style
- **WHEN** a menu item has a destructive variant
- **THEN** that item SHALL render in a destructive visual style (e.g. red text)

#### Scenario: Menu action click does not navigate
- **WHEN** a user clicks a menu item
- **THEN** the menu item's action SHALL fire and navigation SHALL NOT occur
