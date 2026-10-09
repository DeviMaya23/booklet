ALTER TABLE artists
    ADD COLUMN last_used_at TIMESTAMPTZ NULL;

CREATE TABLE artist_links (
    id         UUID        PRIMARY KEY,
    artist_id  UUID        NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    url        TEXT        NOT NULL,
    is_primary BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE artists
    DROP COLUMN artist_link;

UPDATE artists SET last_used_at = NOW();
