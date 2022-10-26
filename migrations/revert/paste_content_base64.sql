-- Revert microbin:paste_content_base64 from pg

BEGIN;

DROP FUNCTION IF EXISTS paste_content_base64;

COMMIT;
