## ADDED Requirements

### Requirement: Struct-tag-driven request validation
The system SHALL validate HTTP request bodies using `go-playground/validator/v10` struct tags. A single validator instance SHALL be registered on the Echo instance at startup and shared across all handlers.

#### Scenario: Validator is registered at startup
- **WHEN** the Echo server initializes
- **THEN** `echo.Validator` is set to an `echoValidator` wrapping a configured `validator.Validate` instance

---

### Requirement: Field names in validation errors use JSON tag names
The validator SHALL be configured with a `RegisterTagNameFunc` that reads the `json` struct tag so that error payloads report field names as they appear in the API (e.g., `"name"` not `"Name"`).

#### Scenario: Validation error field name matches JSON key
- **WHEN** a request fails validation on a field whose `json` tag differs in casing from the Go field name
- **THEN** the error payload reports the `json` tag name, not the Go struct field name

---

### Requirement: Validation errors return structured JSON with HTTP 422
When request validation fails, the system SHALL return HTTP 422 with a JSON body of the form:

```json
{
  "errors": [
    { "field": "<field>", "message": "<human-readable message>" }
  ]
}
```

Bind errors (malformed JSON) SHALL continue to return HTTP 400.

#### Scenario: Validation failure returns 422 with structured errors
- **WHEN** a request body is valid JSON but fails a validation rule
- **THEN** the system returns 422 with `{"errors": [{"field": "...", "message": "..."}]}`

#### Scenario: Malformed JSON still returns 400
- **WHEN** a request body cannot be parsed as JSON
- **THEN** the system returns 400

---

### Requirement: Human-readable validation messages
The system SHALL produce human-readable messages for known validation tags. Unknown tags SHALL fall back to a generic message.

| Tag | Message |
|---|---|
| `required` | `"{field} is required"` |
| `min` | `"{field} must not be empty"` |
| (unknown) | `"{field} failed validation"` |

#### Scenario: Required tag produces readable message
- **WHEN** a `required` field is missing or empty
- **THEN** the error message is `"{field} is required"`

#### Scenario: Min tag on string pointer produces readable message
- **WHEN** a field with `omitempty,min=1` receives an empty string value
- **THEN** the error message is `"{field} must not be empty"`
