-- Revert microbin:appschema from pg

BEGIN;

DROP SCHEMA microbin;

COMMIT;
