CREATE TABLE headers (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    fields jsonb
);

CREATE TABLE blocks (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    fields jsonb
);

-- CREATE TABLE transactions (
--     block_number bigint,
--     block_hash varchar,
--     tx_hash varchar,
--     index integer,
--     fields jsonb
-- );

CREATE TABLE receipts (
    block_number bigint,
    block_hash varchar,
    tx_hash varchar,
    index integer,
    fields jsonb,
    UNIQUE(block_number, index),
    UNIQUE(block_hash, index)
);


CREATE TABLE accounts (
    block_number bigint UNIQUE,
    block_hash varchar UNIQUE,
    address varchar,
    fields jsonb,
    UNIQUE(block_number, address),
    UNIQUE(block_hash, address)
);

CREATE INDEX ON accounts(address);

