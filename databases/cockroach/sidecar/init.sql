ALTER RANGE default CONFIGURE ZONE USING num_replicas = 1;

CREATE DATABASE IF NOT EXISTS kru;

CREATE TABLE IF NOT EXISTS kru.users (
  id STRING PRIMARY KEY,
  username STRING,
  name STRING
);

CREATE USER IF NOT EXISTS kru_admin;
GRANT ALL ON DATABASE kru TO kru_admin;
CREATE USER IF NOT EXISTS kru_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE kru.* TO kru_service;
