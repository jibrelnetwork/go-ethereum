-- SQL in this section is executed when the migration is applied.
-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    DROP TABLE IF EXISTS temp_fix_internal_txs;

    CREATE TABLE temp_fix_internal_txs (
        id SERIAL PRIMARY KEY,
        internal_tx_id INT UNIQUE NOT NULL,
        tx_origin VARCHAR(70) NOT NULL,
        is_processed BOOLEAN NOT NULL
    );

    INSERT INTO temp_fix_internal_txs (internal_tx_id, tx_origin, is_processed) 
    SELECT id, 
    transactions ->> 'from' AS tx_origin, 
    false AS is_processed 
    FROM internal_transactions 
    LEFT JOIN bodies ON internal_transactions.block_hash = bodies.block_hash 
    LEFT JOIN jsonb_array_elements(
        case jsonb_typeof(bodies.fields -> 'Transactions') 
        when 'array' then bodies.fields -> 'Transactions' else '[]' end
    ) transactions ON (transactions ->> 'hash') = internal_transactions.parent_tx_hash 
    WHERE NOT internal_transactions.fields ? 'TxOrigin' and bodies.block_hash IS NOT NULL
    ON CONFLICT DO NOTHING;
END; $$
language plpgsql;

-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS temp_fix_int_txs;
