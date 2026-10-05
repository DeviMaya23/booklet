## REMOVED Requirements

### Requirement: generate_thumbnail job
**Reason**: Images module is being deleted; thumbnail generation for images is no longer needed.
**Migration**: N/A — app is pre-launch. `GenerateThumbnailWorker` and its River job kind `generate_thumbnail` are removed.
