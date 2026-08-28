## Why

`users.id` currently stores the Kinde-issued subject (`kp_abc123`) as a TEXT primary key, fully coupling the app's internal identity to the IdP. An IdP migration would require a database migration across every table with a `user_id` foreign key. Introducing an app-owned UUID as the stable internal key removes that coupling at the database level.

## What Changes

- **BREAKING** `users.id` changes from TEXT (Kinde subject) to UUID (app-generated)
- New column `users.idp_subject TEXT NOT NULL UNIQUE` — stores the Kinde subject, used only at the IdP/bookleaf boundary
- `characters.user_id`, `images.user_id`, `pending_uploads.user_id` FK type changes from TEXT to UUID
- Auth middleware stores the app UUID in context (was: Kinde subject); also stores `idp_subject` as a secondary context value for bookleaf call sites
- `GetOrProvision` resolves users by `idp_subject`, not by `id`
- Bookleaf call sites (`DeleteAccount`, `GetPublicFolders`) continue to use `idp_subject` — no bookleaf API contract change

## Capabilities

### New Capabilities

- `user-identity`: App-owned UUID as the stable internal user key, with `idp_subject` as the IdP lookup bridge

### Modified Capabilities

- `user-auth`: Auth middleware now sets the app UUID (not Kinde subject) as the authenticated user ID in context; additionally stores `idp_subject` in context for cross-system calls

## Impact

- **Database**: All 4 user-related tables; drop-and-recreate is acceptable (no live data)
- **Domain structs**: `User`, `Character`, `Image`, `PendingUpload` — `UserID` field type changes from `string` to `uuid.UUID`
- **Repository layer**: All user/character/image/upload repository method signatures
- **Usecase layer**: All usecase method signatures carrying user identity; `MarkPendingDeletion` takes both `userID uuid.UUID` and `idpSubject string`
- **Auth middleware**: `AuthenticatedUserIDFromContext` return type changes from `string` to `uuid.UUID`; new `AuthenticatedIDPSubjectFromContext` helper added
- **All handlers**: Call sites of `AuthenticatedUserIDFromContext` update to `uuid.UUID`; folder and delete handlers also extract `idpSubject` from context
- **Tests**: Spy/stub method signatures and `domain.User{ID: ...}` literals throughout
