-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS chain_events (
    id bigserial primary key,
    block_number bigint,
    block_hash varchar(70),
    parent_block_hash varchar(70),
    type varchar(20) NOT NULL,
    common_block_number bigint NOT NULL,
    common_block_hash varchar(70) NOT NULL,
    drop_length bigint NOT NULL,
    drop_block_hash varchar(70) NOT NULL,
    add_length bigint NOT NULL,
    add_block_hash varchar(70) NOT NULL,
    node_id varchar(70),
    created_at timestamp default current_timestamp
);

CREATE INDEX IF NOT EXISTS chain_events_block_number_idx ON chain_events(block_number);

CREATE INDEX IF NOT EXISTS chain_events_block_hash_idx ON chain_events(block_hash);

CREATE INDEX IF NOT EXISTS chain_events_parent_block_hash_idx ON chain_events(parent_block_hash);

CREATE INDEX IF NOT EXISTS chain_events_type_idx ON chain_events(type);

CREATE INDEX IF NOT EXISTS chain_events_common_block_number_idx ON chain_events(common_block_number);

CREATE INDEX IF NOT EXISTS chain_events_common_block_hash_idx ON chain_events(common_block_hash);

CREATE INDEX IF NOT EXISTS chain_events_node_id_idx ON chain_events(node_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
