## Why

GET endpoints for characters and images currently return raw R2 storage keys, which are internal identifiers the frontend cannot use to display media. This change replaces raw paths with presigned GET URLs so the frontend can render avatars and images directly from API responses.

## What Changes

- **Characters (GET list, GET by ID):** Replace `avatar_r2_path` with `avatar_url` — a presigned GET URL (1hr TTL). Null when no avatar is set.
- **Images (GET list):** Remove `image_r2_path` entirely. Replace `thumbnail_r2_path` with `thumbnail_url` — a presigned GET URL (1hr TTL). Null when no thumbnail exists.
- **Images (GET by ID):** Replace `image_r2_path` with `image_url` (presigned, 1hr TTL). Replace `thumbnail_r2_path` with `thumbnail_url` (presigned, 1hr TTL). Null when no thumbnail exists.
- **New endpoint — GET /characters/:id/images:** Returns all images tagged with a given character. Each item includes `image_id`, `image_name` (the image's title, nullable), and `thumbnail_url` (presigned, 1hr TTL, nullable).
- **New `Presigner` interface in the handler package:** A shared narrow interface fulfilled by `r2Storage`, injected into handlers that need to presign GET URLs.

## Capabilities

### New Capabilities

- `character-images`: Retrieve all images associated with a character, with presigned thumbnail URLs.

### Modified Capabilities

- `character-management`: Response shape changes — `avatar_r2_path` replaced by `avatar_url`.
- `image-management`: Response shape changes — `image_r2_path` replaced by `image_url` (GET by ID only), `thumbnail_r2_path` replaced by `thumbnail_url`, original file path no longer exposed in list responses.

## Impact

- **API (breaking):** `GET /characters`, `GET /characters/:id`, `GET /images`, `GET /images/:id` response shapes change. Clients relying on raw R2 paths will break.
- **New route:** `GET /characters/:id/images`
- **Handler layer:** New `Presigner` interface; `CharacterHandler` and `ImageHandler` gain a `Presigner` dependency.
- **Repository layer:** New `ListByCharacterID` method on image repository to support the character images endpoint.
- **Wiring:** `main.go` updated to inject `r2Storage` as `Presigner` into the affected handlers.
