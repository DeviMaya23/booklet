CREATE TABLE artists (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    name        TEXT NOT NULL,
    notes       TEXT,
    artist_link TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name)
);

CREATE INDEX idx_artists_user_id ON artists (user_id);

ALTER TABLE images
    DROP COLUMN artist_name,
    DROP COLUMN artist_link,
    ADD COLUMN artist_id UUID REFERENCES artists(id) ON DELETE SET NULL;

ALTER TABLE pending_uploads
    DROP COLUMN artist_name,
    DROP COLUMN artist_link,
    ADD COLUMN artist_id UUID REFERENCES artists(id) ON DELETE SET NULL;
