CREATE TABLE block_headers (
    block_number bigint,
    block_hash varchar,
    fields jsonb
);

CREATE TABLE block_bodies (
    block_number bigint,
    block_hash varchar,
    fields jsonb
);

CREATE TABLE transactions (
    block_number bigint,
    block_hash varchar,
    tx_hash varchar,
    index integer,
    fields jsonb
);

CREATE TABLE receipts (
    block_number bigint,
    block_hash varchar,
    tx_hash varchar,
    index integer,
    fields jsonb
);