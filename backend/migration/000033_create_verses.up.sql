CREATE TABLE verses (
    id           UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id),
    name         TEXT        NOT NULL,
    last_used_at TIMESTAMPTZ NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name)
);

CREATE INDEX idx_verses_user_id ON verses (user_id);
