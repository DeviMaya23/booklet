# Database Schema

All tables carry `created_at` and `updated_at` (TIMESTAMPTZ). Soft-deletable tables also carry `deleted_at`.

```mermaid
erDiagram
    users {
        UUID id PK
        TEXT idp_subject
        TEXT account_state
        TIMESTAMPTZ purged_at
    }
    verses {
        UUID id PK
        UUID user_id FK
        TEXT name
        TIMESTAMPTZ last_used_at
    }
    characters {
        UUID id PK
        UUID user_id FK
        UUID verse_id FK
        TEXT name
        TEXT hero_image_r2_path
        TEXT notes
        BOOLEAN is_public
    }
    character_folders {
        UUID character_id FK
        UUID folder_id
    }
    artists {
        UUID id PK
        UUID user_id FK
        TEXT name
        TEXT notes
        TIMESTAMPTZ last_used_at
    }
    artist_links {
        UUID id PK
        UUID artist_id FK
        TEXT url
        BOOLEAN is_primary
    }
    artpieces {
        UUID id PK
        UUID user_id FK
        UUID artist_id FK
        UUID commission_id FK
        UUID cover_file_id FK
        TEXT title
        TEXT notes
    }
    files {
        UUID id PK
        UUID user_id FK
        UUID artpiece_id FK
        TEXT file_r2_path
        TEXT mime_type
        TEXT thumbnail_r2_path
        TEXT name
        TEXT extension
        TEXT notes
    }
    image_metadata {
        UUID file_id PK FK
        INTEGER width
        INTEGER height
    }
    artpiece_characters {
        UUID artpiece_id FK
        UUID character_id FK
    }
    commissions {
        UUID id PK
        UUID user_id FK
        UUID artist_id FK
        TEXT title
        TEXT status
        NUMERIC price
        BOOLEAN paid
        DATE paid_date
        DATE finish_date
        TIMESTAMPTZ last_contacted_at
        TEXT notes
    }
    commission_characters {
        UUID commission_id FK
        UUID character_id FK
    }
    saved_filters {
        UUID id PK
        UUID user_id FK
        TEXT name
        TEXT thumbnail_r2_path
        JSONB filter_payload
    }
    pending_file_uploads {
        UUID id PK
        UUID user_id FK
        UUID artpiece_id FK
        TEXT r2_key
        TEXT mime_type
        TEXT name
        TEXT extension
        TEXT notes
    }

    users ||--o{ verses : "owns"
    users ||--o{ characters : "owns"
    users ||--o{ artists : "owns"
    users ||--o{ artpieces : "owns"
    users ||--o{ files : "owns"
    users ||--o{ commissions : "owns"
    users ||--o{ saved_filters : "owns"
    users ||--o{ pending_file_uploads : "owns"

    verses ||--o{ characters : "groups"

    characters ||--o{ character_folders : "has"
    characters ||--o{ artpiece_characters : ""
    characters ||--o{ commission_characters : ""

    artists ||--o{ artist_links : "has"
    artists ||--o{ artpieces : "created"
    artists ||--o{ commissions : "contracted"

    artpieces ||--o{ files : "has"
    artpieces ||--o{ artpiece_characters : ""
    artpieces }o--o| commissions : "from"

    files ||--o| image_metadata : "has"

    commissions ||--o{ commission_characters : ""
```

## Notes

- `character_folders.folder_id` is a Bookleaf folder ID (external). No FK — Bookleaf is a separate app.
- `files.artpiece_id` is nullable: unassigned files live in the Inbox.
- `artpieces.cover_file_id` points back into `files` — both sides are nullable (circular, both optional).
- `artpieces.commission_id` is nullable: artpieces can exist without a commission.
- `saved_filters.filter_payload` is a free-form JSONB blob holding the encoded filter state.
- The old `images` / `pending_uploads` tables were dropped in migration `000026`.
