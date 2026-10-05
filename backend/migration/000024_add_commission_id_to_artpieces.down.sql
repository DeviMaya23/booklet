DROP INDEX IF EXISTS idx_artpieces_commission_id;

ALTER TABLE artpieces DROP COLUMN IF EXISTS commission_id;
