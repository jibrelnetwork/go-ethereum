-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE token_holders ADD COLUMN IF NOT EXISTS token_type varchar(15);
ALTER TABLE token_holders ALTER COLUMN token_type SET DEFAULT 'erc20';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
