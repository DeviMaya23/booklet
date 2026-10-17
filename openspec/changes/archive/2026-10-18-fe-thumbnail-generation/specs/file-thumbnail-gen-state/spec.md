# Spec Delta

## REMOVED Requirements

### Requirement: Thumbnail generation state persisted on file record
**Reason**: With client-side thumbnail generation, the `thumbnail_gen_state` lifecycle (pending → done/failed) no longer exists. Thumbnail presence is determined solely by whether `thumbnail_r2_path` is set; the `pending` and `failed` states are not reachable. The column is dropped from the database and the field is removed from all API responses.
**Migration**: Clients should use the presence of `thumbnail_url` in the API response to determine whether a thumbnail exists. Non-image files have no thumbnail; image files have a thumbnail set at upload-complete time.

### Requirement: Thumbnail output format matches source mime type
**Reason**: Format selection logic moves entirely to the client. The backend no longer encodes thumbnails. The format policy (alpha-capable sources produce PNG, others produce JPEG) is now a client responsibility enforced by the frontend upload flow.
**Migration**: Frontend generates thumbnails using OffscreenCanvas: PNG for `image/png`, `image/gif`, `image/webp` sources; JPEG for all others. The thumbnail is PUT to the `thumbnail_upload_url` presigned URL returned by initiate.

### Requirement: Thumbnail generation has a per-job decode timeout
**Reason**: The River worker that this timeout governed is removed. Timeout constraints no longer apply.
**Migration**: None. The constraint is eliminated along with the worker.
