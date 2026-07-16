## Why

Users need to manage images in their booklet — viewing, annotating with artist credits, and associating characters to images. The image entity and its read/update/delete endpoints are a prerequisite for the image upload flow (deferred).

## What Changes

- Introduce the `images` table and `image_characters` join table via a new migration
- Implement four endpoints: list images, get image by ID, update image (partial), delete image
- Image responses embed associated character names (id + name)
- Character associations are managed via full replace on PATCH

## Capabilities

### New Capabilities

- `image-management`: CRUD operations for images (list, get, update, delete). Includes character association management via the `image_characters` join table.

### Modified Capabilities

## Impact

- New files: `domain/image.go`, `usecase/image_usecase.go`, `usecase/image_repository.go`, `repository/image_repository.go`, `handler/image_handler.go`, plus test files for usecase and handler layers
- New DB migration: `images` table + `image_characters` join table
- `main.go`: wire up new repository, usecase, handler, and routes
- No changes to existing character or user modules
- Introduces first use of GORM Many2Many associations in the codebase
