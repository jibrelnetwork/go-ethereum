-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS token_descriptions (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    token varchar(45) NOT NULL,
    total_supply numeric NOT NULL,
    UNIQUE(block_hash, token)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
