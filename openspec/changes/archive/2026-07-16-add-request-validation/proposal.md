## Why

Request validation is currently ad-hoc manual checks scattered in handler code, with no validation at all on the update character path. As new resources are added, this pattern won't scale and leaves gaps.

## What Changes

- Add `github.com/go-playground/validator/v10` as a dependency
- Introduce `handler/validator.go` with an `echoValidator` struct (implements `echo.Validator`), a `RegisterTagNameFunc` that reads `json` tags for field names, and a structured JSON validation error helper
- Wire the validator to the Echo instance at startup
- Replace the manual `if req.Name == ""` check in `CreateCharacter` with a `validate:"required"` struct tag + `c.Validate()`
- Add missing name validation to `UpdateCharacter` via `validate:"omitempty,min=1"` on the `*string` field
- Validation errors return structured JSON: `{"errors": [{"field": "name", "message": "..."}]}`

## Capabilities

### New Capabilities

- `request-validation`: Struct-tag-driven HTTP request validation wired to Echo, with structured JSON error responses

### Modified Capabilities

- `character-management`: Validation rules for create and update request fields are now formally enforced (previously ad-hoc or missing on update)

## Impact

- **New dependency**: `github.com/go-playground/validator/v10`
- **`cmd/server/main.go`**: Wire `e.Validator` after `echo.New()`
- **`internal/handler/validator.go`**: New file — validator wiring + error formatting
- **`internal/handler/character_handler.go`**: Struct tags on request types, `c.Validate()` call, remove manual name check
- **`internal/handler/character_handler_test.go`**: New test cases covering validation error paths
