## 1. Branch & Migration

- [x] 1.1 Checkout `main` and create branch `feat/decouple-idp-identity`
- [x] 1.2 Write migration `000009_decouple_idp_identity.up.sql`: drop `pending_uploads`, `image_characters`, `images`, `character_folders`, `characters`, `users` in that order; recreate `users` with `id UUID PRIMARY KEY` and `idp_subject TEXT NOT NULL UNIQUE`; recreate `characters`, `images`, `pending_uploads` with `user_id UUID NOT NULL REFERENCES users(id)` and all other columns unchanged
- [x] 1.3 Write migration `000009_decouple_idp_identity.down.sql`: drop the same tables and recreate with the old TEXT-based schema

## 2. Domain Structs

- [x] 2.1 Update `domain.User`: change `ID string` to `ID uuid.UUID`; add `IDPSubject string \`gorm:"type:text;not null;uniqueIndex;column:idp_subject"\``
- [x] 2.2 Update `domain.Character`: change `UserID string` to `UserID uuid.UUID \`gorm:"type:uuid;not null;column:user_id"\``
- [x] 2.3 Update `domain.Image`: change `UserID string` to `UserID uuid.UUID \`gorm:"type:uuid;not null;column:user_id"\``
- [x] 2.4 Update `domain.PendingUpload`: change `UserID string` to `UserID uuid.UUID \`gorm:"type:uuid;not null;column:user_id"\``

## 3. Repository Layer

- [x] 3.1 Update `user_repository.GetOrCreate(ctx, idpSubject string)`: generate `uuid.New()` for `ID`, store `idpSubject` in `IDPSubject`; use `ON CONFLICT (idp_subject) DO NOTHING`; then call `GetByIDPSubject` to return the row
- [x] 3.2 Add `user_repository.GetByIDPSubject(ctx, idpSubject string) (*domain.User, error)` — looks up by `idp_subject` column
- [x] 3.3 Update `user_repository.GetByID(ctx, id uuid.UUID)`, `SetPendingDeletion(ctx, id uuid.UUID)`, `DeleteAllUserData(ctx, userID uuid.UUID)` — parameter type change only
- [x] 3.4 Update `UserRepository` interface in `usecase/user_usecase.go` to match: `GetOrCreate(ctx, idpSubject string)`, `GetByIDPSubject(ctx, idpSubject string)`, `GetByID(ctx, id uuid.UUID)`, `SetPendingDeletion(ctx, id uuid.UUID)`, `DeleteAllUserData(ctx, userID uuid.UUID)`
- [x] 3.5 Update `character_repository`: change all `userID string` parameters to `userID uuid.UUID` across `GetByID`, `List`, `Update`, `updateWithFolders`, `Delete`, `GetByIDsAndUserID`
- [x] 3.6 Update `CharacterRepository` and `UploadCharacterRepository` interfaces in `usecase/character_repository.go` to match
- [x] 3.7 Update `image_repository`: change all `userID string` parameters to `userID uuid.UUID` across `GetByID`, `List`, `Update`, `Delete`
- [x] 3.8 Update `ImageRepository` and `UploadImageRepository` interfaces in `usecase/image_repository.go` to match
- [x] 3.9 Update `upload_repository.GetByID(ctx, id uuid.UUID, userID uuid.UUID)`: change `userID` parameter type
- [x] 3.10 Update `UploadRepository` interface in `usecase/upload_repository.go` to match

## 4. Usecase Layer

- [x] 4.1 Update `userUsecase.GetOrProvision(ctx, kindeID string)`: replace `GetByID` + `GetOrCreate` calls with `GetByIDPSubject` + `GetOrCreate(ctx, kindeID)` — resolve by `idp_subject`, not `id`
- [x] 4.2 Update `userUsecase.GetByID` signature to `GetByID(ctx, id uuid.UUID)` and update repo call
- [x] 4.3 Update `userUsecase.MarkPendingDeletion(ctx, userID uuid.UUID, idpSubject string)`: pass `userID` to `SetPendingDeletion` and `idpSubject` to `bookleafClient.DeleteAccount`
- [x] 4.4 Update `userUsecase.PurgeUserData(ctx, userID uuid.UUID)`: parameter type change only
- [x] 4.5 Update `UserUsecase` interface in `handler/user_handler.go` to match new `MarkPendingDeletion` and `PurgeUserData` signatures
- [x] 4.6 Update `UserUsecase` interface in `handler/middleware/auth.go` — `GetOrProvision` return value drives context; no signature change needed but verify it compiles
- [x] 4.7 Update `InitialUploadParams.UserID` from `string` to `uuid.UUID`
- [x] 4.8 Update `uploadUsecase.CompleteUpload(ctx, id uuid.UUID, userID uuid.UUID)`: parameter type change; update `GetByID` and `GetByIDsAndUserID` calls
- [x] 4.9 Update all character usecase methods (`Create`, `GetByID`, `List`, `Update`, `Delete`): change `userID string` to `userID uuid.UUID` in signatures and repo calls
- [x] 4.10 Update all image usecase methods (`GetByID`, `List`, `Update`, `Delete`): change `userID string` to `userID uuid.UUID`
- [x] 4.11 Update `folderUsecase.ListFolders` — signature stays `(ctx, idpSubject string)`; no change needed (the handler will now pass idp_subject from context instead of userID)

## 5. Auth Middleware

- [x] 5.1 Add `AuthenticatedIDPSubjectContextKey ContextKey = "authenticatedIDPSubject"` constant
- [x] 5.2 Add `AuthenticatedIDPSubjectFromContext(c echo.Context) (string, bool)` helper
- [x] 5.3 In `authMiddleware.handle`: replace `c.Set(string(AuthenticatedUserIDContextKey), claims.Subject)` with `c.Set(string(AuthenticatedUserIDContextKey), user.ID)` and add `c.Set(string(AuthenticatedIDPSubjectContextKey), claims.Subject)`
- [x] 5.4 Change `AuthenticatedUserIDFromContext` return type from `(string, bool)` to `(uuid.UUID, bool)`

## 6. Handlers

- [x] 6.1 Update `user_handler.DeleteMe`: extract both `userID` (`uuid.UUID`) and `idpSubject` (`string`) from context; pass both to `MarkPendingDeletion`
- [x] 6.2 Update `user_handler.DeleteUserByID`: `id` path param is now a UUID string — parse with `uuid.Parse` before passing to `PurgeUserData`
- [x] 6.3 Update `character_handler`: change all `userID` variables from `string` to `uuid.UUID` (from updated `AuthenticatedUserIDFromContext`)
- [x] 6.4 Update `image_handler`: same type change as 6.3
- [x] 6.5 Update `upload_handler.InitialUpload`: change `UserID: userID` in `InitialUploadParams` — type now `uuid.UUID`
- [x] 6.6 Update `upload_handler.CompleteUpload`: `userID` from context is now `uuid.UUID`; pass to `CompleteUpload`
- [x] 6.7 Update `folder_handler.ListFolders`: extract `idpSubject` from context using `AuthenticatedIDPSubjectFromContext`; pass it to `ListFolders` instead of `userID`

## 7. Tests

- [x] 7.1 Update `middleware/auth_test.go`: change `domain.User{ID: "sub-1"}` literals to use `uuid.UUID` (e.g. `uuid.New()` or a fixed parse); update `stubUserUsecase.GetOrProvision` return; add scenario asserting `authenticatedIDPSubject` is set in context on active account
- [x] 7.2 Update `usecase/user_usecase_test.go`: update `spyUserRepository` to match new interface (`GetOrCreate(idpSubject string)`, `GetByIDPSubject`, `GetByID(uuid.UUID)`, `SetPendingDeletion(uuid.UUID)`, `DeleteAllUserData(uuid.UUID)`); fix test call sites
- [x] 7.3 Update `usecase/fakes_test.go`: change `fakeCharacterRepository` method signatures from `userID string` to `userID uuid.UUID`; fix `List` comparison (`c.UserID == userID`); update any string-based `uuid.Parse` calls
- [x] 7.4 Update `repository/user_repository_integration_test.go`: replace `"kp_abc123"` and similar Kinde-subject literals with `idpSubject` strings; assert `user.IDPSubject` is set and `user.ID` is a valid UUID; update `seedUser` helper if it sets `ID` directly
- [x] 7.5 Add unit test to `middleware/auth_test.go`: `TestAuth_ActiveAccount_SetsUUIDAndIDPSubjectInContext` — verify `authenticatedUserID` is `user.ID` (uuid.UUID) and `authenticatedIDPSubject` is the JWT subject
- [x] 7.6 Add unit test to `usecase/user_usecase_test.go`: `TestGetOrProvision_NewUser_CreatesWithIDPSubject` — verify provisioning resolves by `idp_subject` and returns a user with a UUID `ID`

## 8. Lint

- [x] 8.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
