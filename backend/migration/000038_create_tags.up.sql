CREATE TABLE tags (
    id           UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id),
    name         TEXT        NOT NULL,
    last_used_at TIMESTAMPTZ NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name),
    CHECK (name = lower(name))
);

CREATE TABLE artpiece_tags (
    artpiece_id UUID NOT NULL REFERENCES artpieces(id) ON DELETE CASCADE,
    tag_id      UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (artpiece_id, tag_id)
);

CREATE INDEX idx_artpiece_tags_tag_id ON artpiece_tags (tag_id);
