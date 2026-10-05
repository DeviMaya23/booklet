# Stale Upload Purge

## Purpose

Defines the background job that periodically finds and hard-deletes pending uploads that were never completed. A pending upload is considered stale once its age exceeds the presign TTL, meaning the client's upload window has closed and the record will never be completed.

---

## Requirements
