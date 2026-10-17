# Design

## Context

See `proposal.md — Why` for motivation.

The thumbnail worker (`internal/worker/generate_file_thumbnail.go`) uses `disintegration/imaging`, which wraps Go's `image.Decode`. Go's decoder registry is extended by blank-importing codec packages. `golang.org/x/image` is already an indirect dependency (pulled in by `disintegration/imaging`); the WebP subpackage just needs a direct blank import to activate.

`domain.File.MimeType` is available on the record fetched at the start of the job, before any image decoding. `imaging.Decode` and `imaging.Fit` are not context-aware — they run synchronously with no cancellation hook.

The existing R2 thumbnail key is always `users/<userID>/thumbnails/<fileID>.jpg`. With format selection, the extension becomes dynamic.

## Goals / Non-Goals

**Goals:**
- Encode thumbnails as PNG when the source is a format that supports alpha; JPEG otherwise
- Register a WebP decoder so `image/webp` source files can be thumbnailed
- Bound the decode+resize step with a hard timeout so a pathological file cannot block a worker goroutine indefinitely

**Non-Goals:**
- Supporting CGO-only formats (AVIF, HEIC, JXL)
- Encoding thumbnails as WebP
- Retroactively re-generating thumbnails for files already in `done` state
- Pixel-scanning to detect actual transparency (as opposed to format capability)

## Decisions

### Format selection uses mime type, not pixel inspection

**Decision**: If `file.MimeType` indicates an alpha-capable format (`image/png`, `image/gif`, `image/webp`), encode the thumbnail as PNG (`image/png`) at a `.png` R2 key. All other image mime types use JPEG encoding at a `.jpg` key, same as today.

**Rationale**: Artists save files as PNG intentionally — format choice carries intent. A pixel scan to check for actual transparency adds O(w×h) latency per job and still doesn't capture intent (an all-opaque PNG may rely on alpha in a future revision). Mime-type selection is O(1), deterministic, and covers the only real-world case: artist PNGs need correct thumbnails.

**Alternative rejected**: Inspect the decoded image's color model via type assertion (e.g. `*image.NRGBA`). `imaging.Decode` normalises all inputs to `*image.NRGBA` regardless of source format, so all decoded images would report the same type — the check would be meaningless.

**Trade-off**: Fully opaque PNGs produce PNG thumbnails, which are modestly larger than equivalent JPEGs at 600px. For a single-user personal library this is acceptable.

### Decode+resize timeout via goroutine and channel

**Decision**: Both `imaging.Decode` and `imaging.Fit` run inside a single goroutine. The worker selects on the goroutine's result channel and the parent context (extended with a 30-second sub-deadline). A timeout is treated the same as any other decode error: River retries, then marks `failed` on exhaustion.

**Rationale**: Neither `imaging.Decode` nor `imaging.Fit` accepts a `context.Context`. The only way to enforce a deadline is to race the operation against a timer. Wrapping both decode and resize together keeps the goroutine count to one and avoids a partial-result state where decoding succeeded but resizing timed out.

**Alternative rejected**: Use a library with context-aware image processing. No pure-Go alternative exists at the required quality; all candidates either require CGO or lack Lanczos resampling.

### WebP via blank import only

**Decision**: Add `_ "golang.org/x/image/webp"` to the worker source file. No version bump is required — `golang.org/x/image` is already an indirect dependency at a version that includes the WebP decoder.

**Rationale**: Go's `image.Decode` uses a registered-decoder pattern. A blank import is sufficient to activate it. Promoting the dependency from indirect to direct requires no code change beyond the import.

## Risks / Trade-offs

- **PNG thumbnails are larger**: A 600×600 PNG thumbnail is roughly 2–4× the size of an equivalent JPEG. Acceptable for the scale of this app; revisit if R2 storage costs become material.
- **Leaked goroutine on timeout**: If the decode goroutine is blocked on a slow `io.Reader` (not CPU), it may not exit immediately when the context fires. The goroutine will eventually unblock when the R2 connection times out or the process exits. The channel is buffered (size 1) so the goroutine can always send and exit without the worker receiving.
- **Existing `.jpg` thumbnails for PNG files remain**: Re-generation is out of scope. Existing files keep their JPEG thumbnail until manually re-triggered or deleted. No migration is needed.

## Migration Plan

No data migration. The change is additive:
- New thumbnails for alpha-capable sources land at `.png` keys alongside existing `.jpg` thumbnails for other files.
- Existing `thumbnail_r2_path` values remain valid.
- Rollback: revert the worker source file and `go.mod`/`go.sum`. Existing thumbnails are unaffected. Any PNG thumbnails generated between deploy and rollback remain in R2 but become unreachable (their file records still point to the `.png` key, which the old worker wouldn't produce again). Acceptable for a personal-use app.
