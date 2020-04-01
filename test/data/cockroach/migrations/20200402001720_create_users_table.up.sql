CREATE TABLE IF NOT EXISTS kiuru.users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  username STRING UNIQUE,
  password STRING,
  name STRING,
  email STRING UNIQUE
);
