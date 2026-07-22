## Why

Characters need to be associated with folders from Bookleaf, an external service, so users can organise their characters into collections. This relationship is managed locally in a join table while the folder data itself lives externally.

## What Changes

- New `character_folders` table linking characters to external folder IDs (composite PK: `character_id`, `folder_id`)
- `Character` gains a `folders` attribute (list of folder IDs) returned on all read endpoints
- `folder_ids` optional field added to create and update character endpoints (UUID shape validation only — no existence check)
- New `GET /folders` endpoint that proxies the Bookleaf internal API and returns the user's available public folders
- New `internal/bookleaf` package as an HTTP client for the Bookleaf internal API
- New `BOOKLEAF_HOST` environment variable for the Bookleaf base URL

## Capabilities

### New Capabilities

- `character-folders`: Association between characters and external Bookleaf folder IDs — create, read, and update folder assignments on characters
- `folder-listing`: Proxy endpoint that fetches a user's available public folders from the Bookleaf internal API

### Modified Capabilities

- `character-management`: Character read endpoints now include a `folders` field; create and update endpoints accept an optional `folder_ids` field

## Impact

- **Database**: New migration for `character_folders` table
- **Config**: New `BOOKLEAF_HOST` env var added to `Config` and `config.go`
- **Domain**: `domain.Character` gains `Folders []CharacterFolder`; new `domain.CharacterFolder` type
- **Repository**: `characterRepository` — Create, GetByID, List, Update, Delete all updated to handle folder rows; repo manages transactions internally
- **Usecase**: `CreateCharacterParams` and `UpdateCharacterParams` gain `FolderIDs *[]uuid.UUID`; new folder usecase proxying the bookleaf client
- **Handler**: Request/response shapes updated; new `GET /folders` endpoint and handler
- **External dependency**: Bookleaf internal API at `GET internal/users/:user_id/public-folders`
- **Tests**: Unit tests on handler and usecase layers; integration tests on character repository
