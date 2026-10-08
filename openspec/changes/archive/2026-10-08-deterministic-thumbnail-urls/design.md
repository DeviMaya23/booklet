## Context

`r2Storage.GeneratePresignedGetURL` wraps `s3.PresignClient.PresignGetObject`, which stamps `time.Now()` as the signing time internally. Every call in the same second produces an identical URL, but every call in a different second differs. Handlers mint thumbnail URLs inline on every list/get response, so the URL changes on every request and the browser cache never hits.

Thumbnail keys are already owner-namespaced: `users/{user_id}/thumbnails/{file_id}.jpg`. No key migration is needed.

All four handlers that serve thumbnail URLs (`file_handler`, `artpiece_handler`, `commission_handler`, `file_upload_handler`) go through the `Presigner` interface. The usecase and worker layers have separate storage interfaces that do not involve presigned GET URLs.

Current TTL for all presigned GET URLs: `PresignGetTTL = 1 hour`.

## Goals / Non-Goals

**Goals:**
- Thumbnail presigned GET URLs are stable for a 24h window — identical URL on every request within the window
- Browser can cache thumbnail responses across navigations and sessions
- No per-instance or shared state required (Cloud Run compatible)
- Full-size file presigned URLs and avatar URLs are unchanged

**Non-Goals:**
- Full-size image deterministic URLs (follow-on)
- Retroactively applying `Cache-Control` to existing thumbnail objects (no migration)
- URL revocation before expiry (accepted trade-off; window capped at 24h)
- Custom thumbnails for non-image files

## Decisions

### Decision: Inject a custom `HTTPPresignerV4` via `s3.PresignOptions.Presigner`

`s3.PresignOptions` accepts a `Presigner HTTPPresignerV4` field. Supplying a thin wrapper that substitutes the caller-provided `signingTime` with a fixed window-start time lets `PresignGetObject` build the request URL normally. The wrapper delegates everything else to `v4.NewSigner()`.

```
windowStart = time.Now().UTC().Truncate(24h)
expiry      = 48h

deterministicPresigner.PresignHTTP(..., signingTime, ...) {
    return inner.PresignHTTP(..., windowStart, ...)
}
```

**Why over building the HTTP request manually with `v4.Signer.PresignHTTP`:** avoids re-implementing R2's URL construction (path-style vs virtual-hosted, query encoding, bucket placement). `PresignGetObject` already handles all of that; we only intercept the time.

**Why over a middleware on `s3.Options.APIOptions`:** the `PresignOptions.Presigner` field is the designed extension point for exactly this; middleware injection is more fragile and undocumented.

### Decision: Add `GenerateDeterministicPresignedGetURL` to the existing `Presigner` handler interface

A single method added to the existing interface. Window size (24h) and expiry (48h) are constants baked into `r2Storage` — callers pass only the key. This keeps the interface surface small and avoids a second interface type and second constructor argument across four handlers.

Character handler gets the method but never calls it — acceptable.

### Decision: Thread Cache-Control through the worker's `fileThumbnailStorageService` interface

Add `cacheControl string` to `PutObject` in `fileThumbnailStorageService`. This interface is unexported and local to the worker package; only one real call site and one fake exist. The worker passes `"private, max-age=172800"` when uploading thumbnails. `r2Storage.PutObject` gains the optional header.

**Why not a separate `PutThumbnailObject` method:** the Cache-Control parameter is general-purpose enough that extending `PutObject` is cleaner than proliferating method variants. The blast radius is one interface, one worker call, one fake.

### Decision: Spike first — verify R2 accepts a fixed signing time

Before touching handlers, a one-off script confirms that calling `PresignGetObject` with a fixed `windowStart` twice (>60s apart) yields byte-identical URLs and that both URLs successfully fetch the object. This gates the rest of the implementation.

## Risks / Trade-offs

- [URL cannot be revoked before expiry] → Acceptable: 48h maximum exposure window; object keys are user-namespaced as defense in depth. Account deletion purges objects.
- [Window rollover at midnight UTC] → Zero user impact: old URLs remain valid for 24–48h after rollover. New API responses return new-window URLs. Browser re-fetches once per device per window, then caches again.
- [Clock skew across Cloud Run instances] → Non-issue at 24h granularity: even seconds of skew produce the same `Truncate(24h)` result.
- [Cache-Control only on new thumbnails] → Accepted. Existing thumbnails get the in-session memory-cache benefit from deterministic URLs; disk-cache persistence requires a re-upload or CopyObject migration (out of scope).
- [R2 `HTTPPresignerV4` compatibility] → Mitigated by the spike task.

## Open Questions

- Does `s3.PresignOptions.Presigner` exist in `aws-sdk-go-v2/service/s3 v1.105.2`? The spike will confirm.
- Should the window duration and expiry be config values rather than constants? Currently: hard-coded constants alongside `PresignGetTTL`. Low priority — they are unlikely to need tuning.
