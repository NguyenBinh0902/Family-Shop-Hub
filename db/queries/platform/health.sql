-- name: Ping :one
-- Used by GET /healthz to prove the API can run a query, not just open a
-- connection.
SELECT 1::int AS ok;
