CREATE TABLE commissions (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    artist_id   UUID REFERENCES artists(id) ON DELETE SET NULL,
    status      TEXT NOT NULL DEFAULT 'waitlist',
    price       NUMERIC,
    paid        BOOLEAN NOT NULL DEFAULT false,
    paid_date   DATE,
    finish_date DATE,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_commissions_user_id ON commissions (user_id);

CREATE TABLE commission_characters (
    commission_id UUID NOT NULL REFERENCES commissions(id) ON DELETE CASCADE,
    character_id  UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    PRIMARY KEY (commission_id, character_id)
);
