## REMOVED Requirements

### Requirement: List images by character
**Reason**: The `GET /characters/:id/images` endpoint is being removed along with the images module. The `imageRepo` dependency in `characterUsecase` is removed. A future change will add a similar endpoint against artpieces/files.
**Migration**: N/A — the frontend never calls this endpoint. Future capability will be a clean addition.
