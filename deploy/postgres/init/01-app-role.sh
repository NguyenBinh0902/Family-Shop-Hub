#!/bin/sh
# Runs once, when the postgres container initializes an empty data directory.
# Creates the application role with LOGIN and a password; migration 0001 then
# grants it privileges. The role is not a superuser and cannot bypass RLS.
set -eu

: "${APP_DB_PASSWORD:?APP_DB_PASSWORD must be set}"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -v app_password="$APP_DB_PASSWORD" <<'SQL'
CREATE ROLE familyshop_app LOGIN PASSWORD :'app_password'
    NOSUPERUSER NOBYPASSRLS NOCREATEDB NOCREATEROLE;
SQL
