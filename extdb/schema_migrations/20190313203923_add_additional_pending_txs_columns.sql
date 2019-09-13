-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS id bigserial primary key;
ALTER TABLE pending_transactions DROP CONSTRAINT IF EXISTS pending_transactions_tx_hash_key;
CREATE INDEX IF NOT EXISTS pending_transactions_tx_hash_idx ON pending_transactions(tx_hash);
ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS timestamp timestamp default current_timestamp;
ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS removed BOOLEAN DEFAULT FALSE;
ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS node_id varchar(70);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
