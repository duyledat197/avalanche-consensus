SELECT 1;
-- CREATE TABLE IF NOT EXISTS nodes (
--   id text PRIMARY KEY,
--   address text NOT NULL,
--   created_at timestamptz,
--   updated_at timestamptz NULL,
--   deleted_at timestamptz
-- );

-- CREATE TABLE IF NOT EXISTS blocks (
--   id text PRIMARY KEY PRIMARY KEY,
--   data INTEGER[],
--   height AUTOINCREMENT,
--   created_at timestamptz,
--   updated_at timestamptz NULL,
--   deleted_at timestamptz
-- );

-- CREATE TABLE IF NOT EXISTS markers (
--   block_id text PRIMARY KEY,
--   created_at timestamptz,
--   updated_at timestamptz NULL,
--   deleted_at timestamptz
-- );