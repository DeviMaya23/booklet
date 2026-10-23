# Tasks

## 1. Handler validation

- [x] 1.1 Add a `validateFilterPayload(p json.RawMessage) *validationErrBody` helper in `saved_filter_handler.go` that returns a field error on `filter_payload` when the raw bytes exceed 4 096 bytes or when `json.Unmarshal` into `map[string]any` fails (non-object value). Verify by running `go build ./...` with no errors.
- [x] 1.2 In `CreateSavedFilter`, call the helper after `c.Validate` succeeds and return 422 if it returns an error. Verify existing handler tests still pass: `go test ./internal/handler/...`.
- [x] 1.3 In `PatchSavedFilter`, call the helper after `c.Validate` succeeds, gated on `req.FilterPayload != nil`. Verify existing handler tests still pass: `go test ./internal/handler/...`.

## 2. Handler tests

- [x] 2.1 Add `TestCreateSavedFilter_NonObjectPayload_Returns422`: POST with `filter_payload` as a JSON array (e.g. `[1,2,3]`). Assert status 422 and that the response body contains a `filter_payload` field error.
- [x] 2.2 Add `TestCreateSavedFilter_OversizedPayload_Returns422`: POST with a `filter_payload` object whose serialized size exceeds 4 096 bytes. Assert status 422 and that the response body contains a `filter_payload` field error.
- [x] 2.3 Add `TestPatchSavedFilter_NonObjectPayload_Returns422`: PATCH with a non-object `filter_payload` (e.g. a JSON string). Assert status 422 and that the response body contains a `filter_payload` field error.
- [x] 2.4 Add `TestPatchSavedFilter_OversizedPayload_Returns422`: PATCH with an oversized `filter_payload`. Assert status 422 and that the response body contains a `filter_payload` field error.
- [x] 2.5 Verify all four new tests pass: `go test ./internal/handler/...`.

## 3. Lint

- [x] 3.1 Run `golangci-lint run ./...` from `backend/` and fix any reported issues. Verify it exits clean.

## Workflow follow-up

- Archive the change once implementation is reviewed and merged.
