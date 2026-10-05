ALTER TABLE artpieces ADD COLUMN commission_id UUID REFERENCES commissions(id) ON DELETE SET NULL;

CREATE INDEX idx_artpieces_commission_id ON artpieces (commission_id);
