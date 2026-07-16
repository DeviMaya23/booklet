## Context

Validation currently lives as manual `if` checks in handler functions. The `CreateCharacter` handler has one manual check (`if req.Name == ""`); `UpdateCharacter` has none, leaving a gap where `{"name": ""}` passes through to the database. As new resources are added, this pattern produces inconsistent validation coverage and inconsistent error responses.

Echo exposes a `Validator` interface (`Validate(i any) error`) that, once registered on the Echo instance, can be invoked per-request via `c.Validate()`. This is the standard integration point.

## Goals / Non-Goals

**Goals:**
- Introduce `go-playground/validator/v10` as the single validation library
- Wire it to Echo so all handlers share one validator instance
- Enforce existing validation rules (create: name required; update: name non-empty if provided) via struct tags
- Return validation errors as structured JSON: `{"errors": [{"field": "name", "message": "..."}]}`

**Non-Goals:**
- Validating fields beyond what the current character spec requires
- Replacing `c.Bind()` — binding and validation remain separate calls
- Custom cross-field or async validators at this stage

## Decisions

### 1. Single `handler/validator.go` for wiring + error formatting

The `echoValidator` struct and the `validationErrorResponse` helper both live in `handler/validator.go`. They're tightly coupled (the formatter is only called when validation fails inside a handler) and neither is large enough to justify its own file now. If the formatter grows (error codes, i18n, etc.), it splits out then.

### 2. `RegisterTagNameFunc` reads `json` struct tags

By default, `validator` uses the Go field name (`Name`). Registering a tag name function:

```go
v.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if name == "-" {
        return ""
    }
    return name
})
```

This makes `FieldError.Field()` return `"name"` instead of `"Name"`, so the JSON error payload matches the API's field naming without post-processing.

**Alternative considered**: lowercase the field name at format time. Rejected — it doesn't handle multi-word fields (e.g. `hero_image_r2_path`) correctly without explicit mapping.

### 3. Human-readable messages via tag switch

The formatter maps known tags to human-readable strings:

| Tag | Message |
|---|---|
| `required` | `"{field} is required"` |
| `min` | `"{field} must not be empty"` |

Unknown tags fall back to `"{field} failed validation"`. This keeps client-facing messages stable and decoupled from library internals.

### 4. `e.Validator` wired in `initEcho`

The validator is registered on the Echo instance in `cmd/server/main.go`'s `initEcho()` function, alongside other Echo-level setup (middleware, CORS). The constructor is `httphandler.NewEchoValidator()` — exported from the handler package, called at startup.

### 5. `c.Validate()` called explicitly after `c.Bind()`

Each handler calls `c.Bind()` then `c.Validate()` as separate steps. This keeps bind errors (malformed JSON → 400) distinct from validation errors (valid JSON, invalid values → 422).

## Risks / Trade-offs

- **Tag coupling**: Struct tags on request types are checked at runtime, not compile time. A typo in a tag silently skips validation. Mitigation: handler unit tests exercise the failing validation path for each tagged field.
- **Message stability**: Human-readable messages are asserted in tests. Adding new tags requires updating the formatter and tests together. Acceptable — the switch is the single source of truth.

## Open Questions

None — scope is fully defined.
