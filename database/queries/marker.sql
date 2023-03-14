-- name: MarkBlock :exec
INSERT INTO
  markers (block_id)
VALUES
  ($1) RETURNING *;

-- name: GetByBlockID :one
SELECT
  *
FROM
  markers
WHERE
  block_id = $1;