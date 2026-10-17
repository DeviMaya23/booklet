# Proposal

## Why

The thumbnail worker has three latent problems: (1) PNG source files with transparency lose their alpha channel because thumbnails are always encoded as JPEG, which has no alpha support — the transparent areas are flattened to black; (2) WebP source files silently fail thumbnail generation because the Go imaging library has no registered WebP decoder, despite `image/webp` being treated as a thumbnailable mime type; (3) a malformed or unusually large source file can block the worker goroutine indefinitely because there is no timeout around the decode/resize step.

## What Changes

- **Backend — thumbnail encoder**: The worker detects whether the decoded source image has an alpha channel. If it does, the thumbnail is encoded as PNG and stored at a `.png` key; otherwise it encodes as JPEG as today. The thumbnail R2 key and content-type are set accordingly.
- **Backend — WebP decoder registration**: A blank import of `golang.org/x/image/webp` is added to the worker package. This registers a pure-Go WebP decoder into `image.Decode`, allowing `imaging.Decode` to process WebP source files with no CGO dependency.
- **Backend — decode timeout**: A `context.WithTimeout` (30 s) wraps the `imaging.Decode` and `imaging.Fit` calls. If the decode or resize exceeds the limit, the job returns an error and River retries or marks it failed on exhaustion — same path as any other decode failure.
- **Backend — go.mod**: `golang.org/x/image` is added (or promoted from indirect) to supply the WebP decoder.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `file-thumbnail-gen-state`: Thumbnail format is now source-dependent (PNG for alpha-channel sources, JPEG otherwise); WebP source files are now decodable; decode/resize step has a hard timeout.

## Impact

- **BE — worker**: `internal/worker/generate_file_thumbnail.go` — format-selection logic, blank WebP import, decode timeout.
- **BE — tests**: `internal/worker/generate_file_thumbnail_test.go` — new scenarios for PNG-with-alpha, WebP source, and decode timeout.
- **BE — go.mod / go.sum**: `golang.org/x/image` added or promoted to direct dependency.
- **No API changes**: `thumbnail_r2_path` remains an opaque string; clients follow it to the CDN regardless of extension. No handler or domain struct changes.
- **No FE changes**.
