## Why

Characters currently have no profile image. The `hero_image_r2_path` column was reserved for this purpose but was never used. This change repurposes it as `avatar_r2_path` and introduces a dedicated two-phase upload flow for character avatars — distinct from image uploads, which carry additional metadata and generate thumbnails.

## What Changes

- Rename `hero_image_r2_path` to `avatar_r2_path` on the `characters` table via migration (column was unused)
- Rename `hero_image_r2_path` to `avatar_r2_path` on the `characters` table via migration (column was unused)
- Add `POST /characters/:id/avatar/init` — initiates an avatar upload, returns a presigned R2 PUT URL
- Add `POST /characters/:id/avatar/:uploadID/complete` — confirms upload; sets `avatar_r2_path` on the character, deletes the previous R2 object if one existed
- Add `DELETE /characters/:id/avatar` — clears the avatar; sets `avatar_r2_path` to null and synchronously deletes the R2 object; idempotent (204 if no avatar exists)
- New `pending_character_avatar_uploads` table to track in-flight uploads (separate from `pending_uploads`, which is image-specific)
- Avatar R2 key format: `:userID/files/:uuid` (no extension — served as-is, no thumbnail generated)
- `avatar_r2_path` is returned in all existing character response shapes
- **BREAKING**: `hero_image_r2_path` is removed from all character request/response fields

## Capabilities

### New Capabilities
- `character-avatar-upload`: Two-phase upload flow for a character's avatar image. Includes init (presign + pending record), complete (character update + old object deletion + pending record cleanup), and delete (clear avatar + synchronous R2 deletion).

### Modified Capabilities
- `character-management`: `hero_image_r2_path` column renamed to `avatar_r2_path`; field appears in all character responses and is accepted in create/update requests.

## Impact

- **DB**: New migration to rename `hero_image_r2_path` → `avatar_r2_path` on `characters`; new `pending_character_avatar_uploads` table
- **Domain**: `Character.HeroImageR2Path` → `Character.AvatarR2Path`
- **Usecase**: `CreateCharacterParams.HeroImageR2Path` → `AvatarR2Path`; `UpdateCharacterParams` same; new `CharacterAvatarUploadUsecase`
- **Handler**: `CharacterHandler` request/response structs updated; new `CharacterAvatarHandler`
- **Storage**: R2 delete called synchronously when a previous avatar exists at complete time
- **Routes**: Three new routes registered under `/characters/:id/avatar`
- **Bruno**: New bruno files for all three new endpoints
