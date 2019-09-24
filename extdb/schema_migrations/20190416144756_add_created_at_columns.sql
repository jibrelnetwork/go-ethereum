-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE headers ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE bodies ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE receipts ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE reorgs ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE chain_splits ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE accounts ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE rewards ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE internal_transactions ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

ALTER TABLE pending_transactions ADD COLUMN IF NOT EXISTS created_at timestamp default current_timestamp;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
