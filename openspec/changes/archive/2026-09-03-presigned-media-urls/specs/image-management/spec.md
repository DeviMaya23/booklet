## MODIFIED Requirements

### Requirement: Image response shape
Every image response SHALL use snake_case field names. The fields included depend on the endpoint:

**GET /images/:id response** SHALL include:
- `id` (string, UUID)
- `image_url` (string) — presigned GET URL for the original file, valid for 1 hour
- `mime_type` (string)
- `title` (string or null)
- `thumbnail_url` (string or null) — presigned GET URL for the thumbnail, valid for 1 hour; null if no thumbnail exists
- `artist_id` (string or null)
- `artist_name` (string or null)
- `notes` (string or null)
- `characters` (array of `{ id: string, name: string }`)
- `created_at` (string, RFC3339)
- `updated_at` (string, RFC3339)

**GET /images response (list)** SHALL include the same fields, except:
- `image_url` is NOT included in list responses — the original file URL is not exposed
- `thumbnail_url` is included (presigned, nullable)

The fields `image_r2_path` and `thumbnail_r2_path` SHALL NOT be present in any image response.

#### Scenario: Response uses snake_case keys
- **WHEN** the system returns an image response
- **THEN** all JSON keys are in snake_case format

#### Scenario: mime_type is present in every image response
- **WHEN** the system returns an image response from any endpoint
- **THEN** the `mime_type` field is present and non-null

#### Scenario: Characters array is present and empty when no associations exist
- **WHEN** an image has no associated characters and is returned by any endpoint
- **THEN** the `characters` field is an empty array `[]`, not null or absent

#### Scenario: artist_id and artist_name are null when no artist is linked
- **WHEN** an image has no artist association and is returned by any endpoint
- **THEN** both `artist_id` and `artist_name` are null

#### Scenario: artist_name is populated from joined artist
- **WHEN** an image has an artist association and is returned by any endpoint
- **THEN** `artist_id` contains the artist's UUID and `artist_name` contains the artist's name

#### Scenario: image_url is a presigned URL in GET by ID
- **WHEN** an authenticated user sends `GET /images/:id` for an image they own
- **THEN** the response includes `image_url` with a presigned HTTPS URL valid for 1 hour

#### Scenario: image_url is absent from list response
- **WHEN** an authenticated user sends `GET /images`
- **THEN** each item in the response does not contain an `image_url` field

#### Scenario: thumbnail_url is a presigned URL when thumbnail exists
- **WHEN** any image endpoint returns an image that has a thumbnail
- **THEN** `thumbnail_url` is a presigned HTTPS URL valid for 1 hour

#### Scenario: thumbnail_url is null when no thumbnail exists
- **WHEN** any image endpoint returns an image with no thumbnail
- **THEN** `thumbnail_url` is null

#### Scenario: raw R2 paths are not exposed in responses
- **WHEN** any image endpoint returns an image
- **THEN** neither `image_r2_path` nor `thumbnail_r2_path` appears in the response
