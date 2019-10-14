-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE index IF NOT EXISTS chain_events_block_number_and_id_idx on chain_events(block_number, id);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
