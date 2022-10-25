-- Revert microbin:paste_type from pg

BEGIN;

DROP TYPE microbin.paste_type;

COMMIT;
