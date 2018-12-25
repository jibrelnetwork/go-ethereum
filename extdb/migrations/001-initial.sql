CREATE TABLE headers (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE bodies (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE pending_transactions (
    tx_hash varchar(70) UNIQUE NOT NULL,
    status varchar,
    fields jsonb
);

CREATE TABLE receipts (
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    fields jsonb
);

CREATE TABLE accounts (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    address varchar(45) NOT NULL,
    fields jsonb,
    UNIQUE(block_hash, address)
);

CREATE TABLE rewards (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) UNIQUE NOT NULL,
    address varchar(45) NOT NULL,
    fields jsonb
);

CREATE TABLE internal_transactions (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    parent_tx_hash varchar(70) NOT NULL,
    index bigint NOT NULL,
    type varchar(20) NOT NULL,
    timestamp bigint NOT NULL,
    fields jsonb
);

CREATE TABLE chain_splits (
    id bigserial primary key,
    common_block_number bigint NOT NULL,
    common_block_hash varchar(70) NOT NULL,
    drop_length bigint NOT NULL,
    drop_block_hash varchar(70) NOT NULL,
    add_length bigint NOT NULL,
    add_block_hash varchar(70) NOT NULL,
    node_id varchar(70)
);

CREATE TABLE reorgs (
    id bigserial primary key,
    block_number bigint NOT NULL,
    block_hash varchar(70) NOT NULL,
    header jsonb,
    reinserted boolean NOT NULL,
    node_id varchar(70)
);


CREATE INDEX ON internal_transactions(block_number);

CREATE UNIQUE INDEX internal_transactions_uniq ON internal_transactions(block_hash, parent_tx_hash, index, type);

CREATE INDEX ON pending_transactions(status);

CREATE INDEX ON internal_transactions(type);

CREATE INDEX ON internal_transactions(block_hash);

CREATE INDEX ON rewards(block_number);

CREATE INDEX ON rewards(block_hash);

CREATE INDEX ON accounts(block_number);

CREATE INDEX ON accounts(block_hash);

CREATE INDEX ON receipts(block_number);

CREATE INDEX ON receipts(block_hash);

CREATE INDEX ON bodies(block_number);

CREATE INDEX ON bodies(block_hash);

CREATE INDEX ON headers(block_number);

CREATE INDEX ON headers(block_hash);

CREATE INDEX ON chain_splits(common_block_number);

CREATE INDEX ON reorgs(block_number);

CREATE INDEX ON reorgs(block_hash);

CREATE INDEX ON reorgs(reinserted);
