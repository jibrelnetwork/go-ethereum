-- +goose Up
-- +goose NO TRANSACTION
-- +goose StatementBegin
CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS internal_transactions_uniq_2 ON internal_transactions(block_hash, parent_tx_hash, index, type, call_depth);
-- +goose StatementEnd

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
