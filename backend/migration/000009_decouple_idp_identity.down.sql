DROP TABLE IF EXISTS pending_uploads;
DROP TABLE IF EXISTS image_characters;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS character_folders;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    account_state TEXT NOT NULL DEFAULT 'active'
        CHECK (account_state IN ('active', 'pending_deletion', 'purged')),
    purged_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE TABLE characters (
    id                UUID PRIMARY KEY,
    user_id           TEXT NOT NULL REFERENCES users(id),
    name              TEXT NOT NULL,
    hero_image_r2_path TEXT,
    biography         TEXT,
    is_public         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_characters_user_id ON characters (user_id);
CREATE INDEX idx_characters_deleted_at ON characters (deleted_at);

CREATE TABLE character_folders (
    character_id UUID NOT NULL REFERENCES characters(id),
    folder_id    UUID NOT NULL,
    PRIMARY KEY (character_id, folder_id)
);

CREATE TABLE images (
    id               UUID PRIMARY KEY,
    user_id          TEXT NOT NULL REFERENCES users(id),
    image_r2_path    TEXT NOT NULL,
    mime_type        TEXT NOT NULL,
    title            TEXT,
    thumbnail_r2_path TEXT,
    artist_name      TEXT,
    artist_link      TEXT,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE image_characters (
    image_id     UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    PRIMARY KEY (image_id, character_id)
);

CREATE INDEX idx_images_user_id ON images (user_id);

CREATE TABLE pending_uploads (
    id            UUID PRIMARY KEY,
    user_id       TEXT NOT NULL REFERENCES users(id),
    r2_key        TEXT NOT NULL,
    mime_type     TEXT NOT NULL,
    title         TEXT,
    artist_name   TEXT,
    artist_link   TEXT,
    notes         TEXT,
    character_ids JSONB NOT NULL DEFAULT '[]',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pending_uploads_user_id ON pending_uploads (user_id);
