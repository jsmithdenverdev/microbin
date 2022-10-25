-- Revert microbin:paste_content_base64 from pg

BEGIN;

DROP FUNCTION IF EXISTS microbin.paste_content_base64;

COMMIT;
