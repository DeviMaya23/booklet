CREATE TABLE images (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    image_r2_path TEXT NOT NULL,
    title TEXT,
    thumbnail_r2_path TEXT,
    artist_name TEXT,
    artist_link TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE image_characters (
    image_id UUID NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    PRIMARY KEY (image_id, character_id)
);

CREATE INDEX idx_images_user_id ON images (user_id);
