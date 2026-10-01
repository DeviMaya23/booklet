## MODIFIED Requirements

### Requirement: Validation errors return structured JSON with HTTP 422
When request validation fails, the system SHALL return HTTP 422 with a JSON body of the form:

```json
{
  "errors": [
    { "field": "<field>", "message": "<human-readable message>" }
  ]
}
```

Bind errors (malformed JSON, or a field whose type cannot be parsed from its JSON value) SHALL return HTTP 400.

UUID-typed fields (`uuid.UUID`, `*uuid.UUID`, `[]uuid.UUID`) SHALL be bound via `encoding.TextUnmarshaler` at bind time. A string value that is not a valid UUID SHALL cause a bind error (HTTP 400), not a validation error (HTTP 422). No field name is included in the error response for bind failures.

#### Scenario: Validation failure returns 422 with structured errors
- **WHEN** a request body is valid JSON but fails a validation rule (e.g., `required`, `min`)
- **THEN** the system returns 422 with `{"errors": [{"field": "...", "message": "..."}]}`

#### Scenario: Malformed JSON returns 400
- **WHEN** a request body cannot be parsed as JSON
- **THEN** the system returns 400

#### Scenario: Invalid UUID value in body returns 400
- **WHEN** a request body is valid JSON but a UUID-typed field contains a string that is not a valid UUID
- **THEN** the system returns 400 with `{"message": "invalid request body"}`
- **AND** no field name is included in the error response
