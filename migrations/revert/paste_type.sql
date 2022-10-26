-- Revert microbin:paste_type from pg

BEGIN;

DROP TYPE paste_type;

COMMIT;
