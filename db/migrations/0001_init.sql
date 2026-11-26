-- Base setup: the non-superuser application role and the RLS helper.
--
-- Migrations run as the database owner. The API connects as familyshop_app,
-- which is neither superuser nor BYPASSRLS, so Row-Level Security policies
-- really apply to it (NFR-05, ADR-0003). In docker-compose the role is
-- created with LOGIN and a password by deploy/postgres/init; this migration
-- only creates it if missing (e.g. a fresh test database) and hardens it.

-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'familyshop_app') THEN
        CREATE ROLE familyshop_app NOLOGIN;
    END IF;
END
$$;
-- +goose StatementEnd

ALTER ROLE familyshop_app NOSUPERUSER NOBYPASSRLS NOCREATEDB NOCREATEROLE;

GRANT USAGE ON SCHEMA public TO familyshop_app;

-- Every table and sequence the owner creates from now on is usable by the
-- application role without a GRANT in each migration. Tables that must be
-- append-only revoke UPDATE/DELETE explicitly in their own migration.
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO familyshop_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO familyshop_app;

-- app_current_household returns the household of the current transaction, set
-- by platform/db with set_config('app.household_id', ..., true). It returns
-- NULL outside a tenant transaction, so policies comparing against it match
-- no rows. All RLS policies use it:
--   USING (household_id = app_current_household())
CREATE FUNCTION app_current_household() RETURNS uuid
    LANGUAGE sql STABLE
    AS $$ SELECT nullif(current_setting('app.household_id', true), '')::uuid $$;

-- +goose Down
DROP FUNCTION app_current_household();

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE USAGE, SELECT ON SEQUENCES FROM familyshop_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM familyshop_app;

REVOKE USAGE ON SCHEMA public FROM familyshop_app;
-- The role itself is cluster-wide and may be in use by other databases, so it
-- is not dropped.
