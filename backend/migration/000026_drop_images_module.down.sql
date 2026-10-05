CREATE TABLE images (
    id                UUID PRIMARY KEY,
    user_id           UUID NOT NULL REFERENCES users(id),
    image_r2_path     TEXT NOT NULL,
    mime_type         TEXT NOT NULL,
    title             TEXT,
    thumbnail_r2_path TEXT,
    artist_id         UUID REFERENCES artists(id) ON DELETE SET NULL,
    notes             TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_images_user_id ON images (user_id);

CREATE TABLE image_characters (
    image_id     UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    PRIMARY KEY (image_id, character_id)
);

CREATE TABLE pending_uploads (
    id           UUID PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES users(id),
    r2_key       TEXT NOT NULL,
    mime_type    TEXT NOT NULL,
    title        TEXT,
    artist_id    UUID REFERENCES artists(id) ON DELETE SET NULL,
    notes        TEXT,
    character_ids JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pending_uploads_user_id ON pending_uploads (user_id);
