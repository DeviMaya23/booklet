## 1. Dependency

- [x] 1.1 Add `github.com/go-playground/validator/v10` to `go.mod` via `go get`

## 2. Validator Infrastructure

- [x] 2.1 Create `internal/handler/validator.go` with `echoValidator` struct implementing `echo.Validator`, configured with `RegisterTagNameFunc` reading `json` struct tags
- [x] 2.2 Add `validationErrResponse` helper in `validator.go` that converts `validator.ValidationErrors` to `{"errors": [{"field": "...", "message": "..."}]}` with human-readable messages for `required` and `min` tags
- [x] 2.3 Wire `e.Validator = httphandler.NewEchoValidator()` in `initEcho` in `cmd/server/main.go`

## 3. Character Handler

- [x] 3.1 Add `validate:"required"` tag to `Name` field on `createCharacterRequest`; add `validate:"omitempty,min=1"` tag to `Name *string` field on `updateCharacterRequest`
- [x] 3.2 In `CreateCharacter`: replace the manual `if req.Name == ""` check with `c.Validate(&req)`, returning 422 with structured error on failure
- [x] 3.3 In `UpdateCharacter`: add `c.Validate(&req)` after `c.Bind`, returning 422 with structured error on failure

## 4. Tests

- [x] 4.1 In `character_handler_test.go`, add test case for `POST /characters` with empty name → 422 with `{"errors": [{"field": "name", "message": "name is required"}]}`
- [x] 4.2 In `character_handler_test.go`, add test case for `PATCH /characters/:id` with `"name": ""` → 422 with `{"errors": [{"field": "name", "message": "name must not be empty"}]}`
- [x] 4.3 In `character_handler_test.go`, add test case for `PATCH /characters/:id` without name field → 200 (name absent is valid)

## 5. Lint

- [x] 5.1 Run `golangci-lint run ./...` from `backend/` and fix any issues
