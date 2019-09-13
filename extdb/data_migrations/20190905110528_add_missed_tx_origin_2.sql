-- SQL in this section is executed when the migration is applied.
-- +goose Up
-- +goose StatementBegin

DO $$
DECLARE
   current_id integer = 0;
   max_id integer = 0;
   batch_size integer = 10000;
BEGIN
   SELECT MIN(internal_tx_id) INTO current_id
   FROM temp_fix_internal_txs
   WHERE NOT is_processed;

   SELECT MAX(internal_tx_id) INTO max_id
   FROM temp_fix_internal_txs
   WHERE NOT is_processed;

   WHILE current_id < max_id
   LOOP
        RAISE LOG 'fix_internal_txs: Add TxOrigin for ID: %', current_id;

        UPDATE internal_transactions
        SET fields = fields || jsonb_build_object('TxOrigin', tx_origin)
        FROM (
            SELECT internal_tx_id, tx_origin
            FROM temp_fix_internal_txs
            WHERE NOT is_processed AND internal_tx_id BETWEEN current_id
            AND (current_id + batch_size)
        ) AS subquery
        WHERE internal_transactions.id=subquery.internal_tx_id;

        UPDATE temp_fix_internal_txs
        SET is_processed = true
        WHERE internal_tx_id BETWEEN current_id AND (current_id + batch_size);

        current_id := current_id + batch_size;
   END LOOP; 
   DROP TABLE IF EXISTS temp_fix_internal_txs;
END; $$
language plpgsql;

-- +goose StatementEnd
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
