-- Revert microbin:paste from pg

BEGIN;

DROP TABLE IF EXISTS microbin.paste;
COMMIT;
