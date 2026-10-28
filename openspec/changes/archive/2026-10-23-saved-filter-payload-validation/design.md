# Design

## Context

See proposal.md — Why.

`filter_payload` is stored as `jsonb` in Postgres and typed as `json.RawMessage` in both request structs. Echo's `c.Bind` unmarshals the whole request body as JSON, so a structurally malformed body (broken JSON syntax) already fails at bind time with HTTP 400. What is not covered: a syntactically valid JSON value that is not an object (string, array, null literal), and a payload of unbounded size.

The existing validator (`go-playground/validator`) uses struct tags. `json.RawMessage` is a `[]byte`, and the standard `max` tag operates on string length or slice element count — not byte length of the raw JSON. A custom type func or registration would be needed to use tags for this.

## Goals / Non-Goals

**Goals:**
- Reject a `filter_payload` that is not a JSON object (returns 422 with a `filter_payload` field error)
- Reject a `filter_payload` whose raw byte length exceeds 4 096 bytes (returns 422 with a `filter_payload` field error)
- Apply both checks on create and on patch (only when `filter_payload` is present)

**Non-Goals:**
- Validating the internal structure of the payload (specific keys, value types) — the FE always sends the known shape; deep validation adds coupling with no real benefit
- Adding a general-purpose `json_object` validator tag to the shared validator

## Decisions

### Inline handler validation after Bind, not a custom validator tag

**Decision**: add a small helper `validateFilterPayload(p json.RawMessage) error` and call it in `CreateSavedFilter` and `PatchSavedFilter` after `c.Bind` / `c.Validate` succeed.

**Why**: the project already uses this pattern for cases `c.Bind` doesn't cover (UUID path params, empty-name patch check). Registering a custom type func for `json.RawMessage` would touch the shared validator and affect every handler; the check is specific to this one field.

**Alternative considered**: `validate:"max=4096"` tag on a custom type — rejected because `max` on `[]byte` counts elements, not raw JSON bytes, and a custom type func adds shared-validator complexity for a single field.

**Error format**: return the standard `validationErrBody` shape (`{"errors":[{"field":"filter_payload","message":"..."}]}`) with HTTP 422 — consistent with all other field errors in this codebase.

## Risks / Trade-offs

- The 4 096 byte cap is generous relative to the current schema and leaves ample room for future filter fields. If the schema ever grows to structured sub-objects, the cap can be raised in one place (the handler constant).
- `json.Unmarshal` into `map[string]any` to confirm object shape is a small allocation per request on the two mutating endpoints only — negligible.
