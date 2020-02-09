#!/bin/sh

exec service \
  --port "$PORT" \
  --gateway-port "$GATEWAY_PORT" \
  --crdb-host "$CRDB_HOST" \
  --crdb-user "$CRDB_USER" \
  --crdb-dbname "$CRDB_DBNAME" \
  --redis-host "$REDIS_HOST" \
  --jwt-secret "$JWT_SECRET"
