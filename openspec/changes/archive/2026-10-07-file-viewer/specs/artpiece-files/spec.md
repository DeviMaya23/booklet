## ADDED Requirements

### Requirement: Artpiece file response includes name and image dimensions
The artpiece detail response SHALL include `name`, `width`, and `height` fields for each attached file.

#### Scenario: File response with image dimensions
- **WHEN** `GET /artpieces/:id` is called and a file has stored `ImageMetadata`
- **THEN** each file entry in the response `files` array SHALL include `name` (string or null), `width` (integer or null), and `height` (integer or null)

#### Scenario: File response without image dimensions
- **WHEN** `GET /artpieces/:id` is called and a file has no `ImageMetadata` (e.g. PSD, PDF)
- **THEN** the file entry SHALL include `name` (string or null), and `width` and `height` SHALL both be null
