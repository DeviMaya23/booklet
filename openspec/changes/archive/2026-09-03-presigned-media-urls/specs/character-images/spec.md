## ADDED Requirements

### Requirement: List images by character
An authenticated user SHALL be able to retrieve all images associated with a given character. The character is identified by its UUID path parameter. The endpoint SHALL return only images owned by the requesting user.

The response is an array of image summary objects. Each object SHALL include:
- `image_id` (string, UUID)
- `image_name` (string or null) — the image's `title` field; null if no title has been set
- `thumbnail_url` (string or null) — presigned GET URL for the thumbnail, valid for 1 hour; null if no thumbnail exists

If the character does not exist, belongs to another user, or has no associated images owned by the user, the endpoint SHALL return an empty array.

#### Scenario: Returns images tagged with the character
- **WHEN** an authenticated user sends `GET /characters/:id/images` for a character they own that has associated images
- **THEN** the system returns 200 with an array of image summary objects for each associated image

#### Scenario: Returns empty array when no images are tagged
- **WHEN** an authenticated user sends `GET /characters/:id/images` for a character they own that has no associated images
- **THEN** the system returns 200 with an empty array

#### Scenario: Returns empty array for unknown or unowned character
- **WHEN** an authenticated user sends `GET /characters/:id/images` for a character that does not exist or belongs to another user
- **THEN** the system returns 200 with an empty array

#### Scenario: thumbnail_url is presigned when thumbnail exists
- **WHEN** the endpoint returns an image that has a thumbnail
- **THEN** `thumbnail_url` is a presigned HTTPS URL valid for 1 hour

#### Scenario: thumbnail_url is null when no thumbnail exists
- **WHEN** the endpoint returns an image that has no thumbnail
- **THEN** `thumbnail_url` is null

#### Scenario: image_name is null when title is not set
- **WHEN** the endpoint returns an image with no title
- **THEN** `image_name` is null

#### Scenario: Invalid UUID path param
- **WHEN** an authenticated user sends `GET /characters/:id/images` with a non-UUID character id
- **THEN** the system returns 400
