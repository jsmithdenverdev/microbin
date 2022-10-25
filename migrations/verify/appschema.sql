-- Verify microbin:appschema on pg

BEGIN;

SELECT pg_catalog.has_schema_privilege('microbin', 'usage');

ROLLBACK;
