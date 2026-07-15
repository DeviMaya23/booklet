CREATE TABLE characters (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    hero_image_r2_path TEXT,
    biography TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_characters_user_id ON characters (user_id);
CREATE INDEX idx_characters_deleted_at ON characters (deleted_at);
