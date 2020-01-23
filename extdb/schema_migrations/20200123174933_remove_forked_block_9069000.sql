-- +goose Up
-- SQL in this section is executed when the migration is applied.
DELETE FROM chain_events WHERE block_number=9069000 AND block_hash='0x072cf1df374159c5f23087750d8a2f3201542da196939ce446ff2c5c390fe5f6';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
