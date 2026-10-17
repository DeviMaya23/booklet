# Spec Delta

## Purpose

Defines the FileThumb component — a category-aware file thumbnail that shows the best available representation for a file: a generated image thumbnail, a loading skeleton while one is pending, or a category tile (icon, optional color tint, and extension label) when no thumbnail exists or is applicable.

## ADDED Requirements

### Requirement: Category tile fallback when no thumbnail is available
When a file has no `thumbnail_url` and `thumbnail_gen_state` is not `pending`, the system SHALL render a category tile in place of the thumbnail. The tile SHALL display a category-appropriate icon, a color-tinted background for the four high-signal categories (Video, Audio, Layered, 3D), and the file's extension as an uppercase pill label when `extension` is non-null.

#### Scenario: Non-image file with known extension shows category icon and label
- **WHEN** a file tile renders with no `thumbnail_url`, `thumbnail_gen_state` of `not_applicable`, and `extension` of `"clip"`
- **THEN** the tile displays the Layered category icon with a terracotta tint background and a "CLIP" label pill

#### Scenario: Image file with failed thumbnail shows Image category tile
- **WHEN** a file tile renders with no `thumbnail_url`, `thumbnail_gen_state` of `"failed"`, and `extension` of `"png"`
- **THEN** the tile displays the Image category icon on a neutral background with a "PNG" label pill

#### Scenario: File with unknown extension shows fallback icon and no label
- **WHEN** a file tile renders with no `thumbnail_url`, `thumbnail_gen_state` of `not_applicable`, and `extension` of `null`
- **THEN** the tile displays the generic fallback icon on a neutral background with no label pill

### Requirement: Loading indicator shown while thumbnail generation is pending
When a file has no `thumbnail_url` and `thumbnail_gen_state` is `"pending"`, the system SHALL render an animated loading spinner in place of the category tile.

#### Scenario: Pending thumbnail shows loading spinner
- **WHEN** a file tile renders with no `thumbnail_url` and `thumbnail_gen_state` of `"pending"`
- **THEN** the tile displays an animated loading spinner, not a category icon

### Requirement: Generated thumbnail takes precedence over category tile
When a file has a `thumbnail_url`, the system SHALL render the thumbnail image regardless of file category.

#### Scenario: Image file with thumbnail shows the thumbnail
- **WHEN** a file tile renders with a non-null `thumbnail_url`
- **THEN** the tile displays the thumbnail image, not a category icon or skeleton

### Requirement: Category tile is theme-aware
The category tile's tint colors SHALL be defined as CSS custom properties and SHALL produce a visually distinguishable tile in every supported theme, including the dark theme.

#### Scenario: Tint colors are distinct from page background in dark theme
- **WHEN** the active theme is `sunless` and a Video-category file is displayed without a thumbnail
- **THEN** the tile background color is visibly different from the page background color

### Requirement: FileThumb is used at all file tile rendering sites
The system SHALL use a single shared FileThumb component for the thumbnail fallback in all file tile contexts: the inbox grid, the artpiece detail file grid, the artpiece files input, and the file tile edit overlay.

#### Scenario: Inbox grid renders category tile via FileThumb
- **WHEN** the inbox grid displays a file with no thumbnail and a known extension
- **THEN** the category tile is rendered (same behavior as artpiece detail grid for the same file)
