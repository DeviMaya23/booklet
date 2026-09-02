CREATE TABLE pending_character_avatar_uploads (
    id           uuid        PRIMARY KEY,
    user_id      uuid        NOT NULL REFERENCES users(id),
    character_id uuid        NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    r2_key       text        NOT NULL,
    mime_type    text        NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT NOW()
);
