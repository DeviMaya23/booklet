## Purpose

This capability covers attaching and detaching files from artpieces, as well as explicitly setting an artpiece's cover file. It includes auto-selection logic for cover assignment when files are attached or detached.

## Requirements

### Requirement: Attach file to artpiece
The system SHALL allow an authenticated user to attach an existing file to an artpiece, both owned by the same user. After attachment, cover auto-selection runs if the artpiece has no cover.

#### Scenario: Successful attach — artpiece has no cover, file is an image
- **WHEN** a user attaches an image-type file to an artpiece that currently has no cover
- **THEN** the file is attached (artpiece_id set) and the artpiece's cover_file_id is set to that file

#### Scenario: Successful attach — artpiece has no cover, file is not an image, no other images exist
- **WHEN** a user attaches a non-image file to a coverless artpiece that contains no other image files
- **THEN** the file is attached and set as the cover

#### Scenario: Successful attach — artpiece already has a cover
- **WHEN** a user attaches a file to an artpiece that already has a cover
- **THEN** the file is attached and the existing cover is left unchanged

#### Scenario: File does not belong to user
- **WHEN** a user attaches a file that belongs to another user
- **THEN** the system returns HTTP 422

#### Scenario: Artpiece does not belong to user
- **WHEN** a user attaches a file to an artpiece that belongs to another user
- **THEN** the system returns HTTP 404

#### Scenario: File already attached to a different artpiece
- **WHEN** a user attaches a file that already has a different artpiece_id
- **THEN** the system returns HTTP 422

---

### Requirement: Detach file from artpiece
The system SHALL allow an authenticated user to detach a file from an artpiece. Detaching sets the file's artpiece_id to NULL. If the detached file was the cover, cover reassignment runs.

#### Scenario: Successful detach — file was not the cover
- **WHEN** a user detaches a file that is not the artpiece's cover
- **THEN** the file's artpiece_id is set to NULL and the artpiece cover is unchanged

#### Scenario: Cover reassignment — another image file exists
- **WHEN** a user detaches the cover file and another image-type file is still attached
- **THEN** the detached file's artpiece_id is set to NULL and the artpiece cover is reassigned to another image file

#### Scenario: Cover reassignment — no images remain, another file exists
- **WHEN** a user detaches the cover file and no image-type files remain but other files do
- **THEN** the artpiece cover is set to the first remaining file

#### Scenario: Cover cleared — no files remain
- **WHEN** a user detaches the only file from the artpiece
- **THEN** the file's artpiece_id is set to NULL and the artpiece cover is set to NULL

#### Scenario: File not attached to this artpiece
- **WHEN** a user attempts to detach a file that is not attached to the specified artpiece
- **THEN** the system returns HTTP 422

#### Scenario: Artpiece does not belong to user
- **WHEN** a user detaches a file from an artpiece that belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Set cover
The system SHALL allow an authenticated user to explicitly set any file attached to an artpiece as its cover.

#### Scenario: Successful cover set
- **WHEN** a user sets a file as the cover and that file's artpiece_id equals the artpiece's id
- **THEN** the artpiece's cover_file_id is updated to that file and the system returns the updated artpiece

#### Scenario: File not attached to this artpiece
- **WHEN** a user sets a cover using a file whose artpiece_id does not match the artpiece
- **THEN** the system returns HTTP 422

#### Scenario: File does not belong to user
- **WHEN** a user sets a cover using a file that belongs to another user
- **THEN** the system returns HTTP 422

#### Scenario: Artpiece does not belong to user
- **WHEN** a user sets a cover on an artpiece that belongs to another user
- **THEN** the system returns HTTP 404

---

### Requirement: Artpiece create accepts optional initial file IDs
The system SHALL allow `POST /artpieces` to accept an optional `file_ids` array. When provided, the files are attached to the new artpiece and cover assignment rules are applied. An absent or empty array SHALL be equivalent to a plain create with no files.

#### Scenario: Create with no file_ids
- **WHEN** a user creates an artpiece without a `file_ids` field (or with an empty array)
- **THEN** the system creates the artpiece with no files attached, returning 201

#### Scenario: Create with valid file IDs
- **WHEN** a user creates an artpiece with a non-empty `file_ids` array of files they own
- **THEN** the system creates the artpiece, attaches all specified files, assigns the cover to the first image file in the set (or the first file if no images exist), and returns 201

#### Scenario: Create with a file ID not owned by the user
- **WHEN** a user creates an artpiece and one or more `file_ids` entries do not exist or belong to another user
- **THEN** the system returns 422 and does not create the artpiece

#### Scenario: Create with a file already attached to another artpiece
- **WHEN** a user creates an artpiece and one or more `file_ids` entries are already attached to a different artpiece
- **THEN** the system returns 422 and does not create the artpiece

---

### Requirement: Bulk full-replace artpiece file set
The system SHALL provide `PUT /artpieces/:id/files` that accepts the complete desired set of file IDs and reconciles it in a single transaction: files newly included are attached, files no longer included are detached, and unchanged files are left alone. Cover assignment rules are applied to the resulting set.

#### Scenario: Full replace with a new set
- **WHEN** a user sends `PUT /artpieces/:id/files` with a `file_ids` array
- **THEN** the system attaches files in the new set not currently attached, detaches files currently attached but absent from the new set, and returns 200 with the updated artpiece

#### Scenario: Empty new set
- **WHEN** a user sends `PUT /artpieces/:id/files` with an empty `file_ids` array
- **THEN** all files are detached from the artpiece, the cover is cleared, and the system returns 200

#### Scenario: Current cover remains in the new set
- **WHEN** the artpiece's current cover file is included in the new `file_ids`
- **THEN** the cover is unchanged

#### Scenario: Current cover removed from set
- **WHEN** the artpiece's current cover file is not in the new `file_ids`
- **THEN** the system reassigns the cover to the first image in the remaining set, or the first file if no images, or clears it if the set is empty

#### Scenario: Artpiece had no cover and new set contains images
- **WHEN** the artpiece has no cover and the new set contains at least one image file
- **THEN** the system sets the cover to the first image file in the new set

#### Scenario: File in new set belongs to another user
- **WHEN** one or more `file_ids` entries do not exist or belong to another user
- **THEN** the system returns 422 and makes no changes

#### Scenario: File in new set already attached to a different artpiece
- **WHEN** one or more `file_ids` entries are already attached to a different artpiece
- **THEN** the system returns 422 and makes no changes

#### Scenario: Artpiece not found
- **WHEN** a user sends `PUT /artpieces/:id/files` for an artpiece that does not exist or belongs to another user
- **THEN** the system returns 404
