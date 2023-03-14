-- name: GetLatestBlock :many
SELECT
  *
FROM
  blocks
ORDER BY
  created_at
LIMIT
  $1;