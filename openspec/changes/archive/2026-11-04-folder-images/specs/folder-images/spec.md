# Spec Delta

## Purpose

Defines the rules for fetching the thumbnail image list of a Bookleaf public folder via booklet, including the Bookleaf client method, the proxy endpoint, and the frontend query hook.

## ADDED Requirements

### Requirement: Fetch folder image list via proxy endpoint
An authenticated user SHALL be able to retrieve the images of one of their public Bookleaf folders via `GET /folders/:folderID/images`. The system SHALL call the Bookleaf internal API and return the image list without transformation.

The Bookleaf API called is: `GET <BOOKLEAF_HOST>/internal/users/:userID/folders/:folderID/images`

The user ID sent to Bookleaf SHALL be the authenticated user's IDP subject from the JWT token.

The response shape returned to the client is:
```json
{
  "images": [
    {
      "image_id": "<string>",
      "thumbnail_url": "<presigned URL>"
    }
  ]
}
```

An empty folder SHALL return `{ "images": [] }` — not a 404.

#### Scenario: Successful image list for a non-empty folder
- **WHEN** an authenticated user sends `GET /folders/:folderID/images` for a folder that has images
- **THEN** the system returns 200 with the list of images, each containing `image_id` and `thumbnail_url`

#### Scenario: Successful image list for an empty folder
- **WHEN** an authenticated user sends `GET /folders/:folderID/images` for a folder that has no images
- **THEN** the system returns 200 with `{ "images": [] }`

#### Scenario: Folder not found or not owned by the user
- **WHEN** Bookleaf returns 404 for the given folder ID and user ID combination
- **THEN** the system returns 404 to the caller

#### Scenario: Bookleaf API is unavailable
- **WHEN** the Bookleaf API returns an error or is unreachable
- **THEN** the system returns 500

---

### Requirement: Bookleaf client method for folder images
The Bookleaf HTTP client SHALL expose a `GetFolderImages(ctx, userID, folderID)` method that calls `GET <BOOKLEAF_HOST>/internal/users/:userID/folders/:folderID/images` and returns a typed result. The same `X-Bookleaf-Internal-Secret` header pattern used by existing client methods SHALL be applied.

The response is decoded into a struct containing a list of images, each with `image_id` and `thumbnail_url` string fields.

#### Scenario: Successful response is decoded
- **WHEN** Bookleaf returns 200 with a valid image list
- **THEN** the client returns the decoded image list with no error

#### Scenario: Non-200 response from Bookleaf is surfaced as an error
- **WHEN** Bookleaf returns a non-200 status
- **THEN** the client returns an error that includes the HTTP status code

---

### Requirement: Frontend query hook for folder images
The frontend SHALL provide a `useGetFolderImages(folderID)` React Query hook that calls `GET /folders/:folderID/images` and returns the image list. The hook SHALL only fetch when `folderID` is a non-empty string. The hook SHALL not retry on failure.

#### Scenario: Hook fetches when folderID is provided
- **WHEN** `useGetFolderImages` is called with a non-empty `folderID`
- **THEN** the hook issues a request to `GET /folders/:folderID/images` and returns the image list on success

#### Scenario: Hook does not fetch when folderID is empty
- **WHEN** `useGetFolderImages` is called with an empty or undefined `folderID`
- **THEN** no request is issued
