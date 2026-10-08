CREATE TABLE saved_filters (
    id                uuid        PRIMARY KEY,
    user_id           uuid        NOT NULL REFERENCES users(id),
    name              text        NOT NULL,
    thumbnail_r2_path text,
    filter_payload    jsonb       NOT NULL,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);
