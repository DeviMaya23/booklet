# Proposal

## Why

The `filter_payload` field on saved filter create and patch endpoints accepts any JSON value without checking its shape or size. A client could send a non-object value (string, array, null) or an arbitrarily large blob; both pass through to Postgres silently.

## What Changes

- On **create** (`POST /saved_filters`): after bind, validate that `filter_payload` is a JSON object and does not exceed 4 KB. Return HTTP 422 with a structured field error if either check fails.
- On **patch** (`PATCH /saved_filters/:id`): same two checks, applied only when `filter_payload` is present in the request.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `saved-filters`: add requirements that `filter_payload` must be a JSON object and must be ≤ 4 KB on create and update.

## Impact

- `backend/internal/handler/saved_filter_handler.go`: inline validation after `c.Bind` in `CreateSavedFilter` and `PatchSavedFilter`
- `backend/internal/handler/saved_filter_handler_test.go`: new handler test cases for shape rejection and size rejection
