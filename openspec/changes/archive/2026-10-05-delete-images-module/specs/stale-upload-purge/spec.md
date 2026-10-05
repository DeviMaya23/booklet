## REMOVED Requirements

### Requirement: Stale identification
**Reason**: The `pending_uploads` table is being dropped with the images module. Stale upload purge for images is no longer needed. The analogous purge for `pending_file_uploads` (file-upload module) is unaffected and remains.
**Migration**: N/A — app is pre-launch. `PurgeExpiredUploadsWorker` and its periodic job are removed from `main.go`.

### Requirement: Deletion order
**Reason**: See above.
**Migration**: N/A.

### Requirement: Cleanup logging
**Reason**: See above.
**Migration**: N/A.

### Requirement: Retry on failure
**Reason**: See above.
**Migration**: N/A.
