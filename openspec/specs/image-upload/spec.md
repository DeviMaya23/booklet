# Image Upload

## Purpose

Defines the two-phase upload flow for images. A user first initiates an upload to receive a presigned R2 PUT URL (InitialUpload), then confirms completion after uploading the file directly to R2 (CompleteUpload). This decouples metadata persistence from file transfer.

---

## Requirements
