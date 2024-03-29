#!/bin/sh

exec service \
  --port "$PORT" \
  --gateway-port "$GATEWAY_PORT" \
  --crdb-host "$CRDB_HOST" \
  --crdb-user "$CRDB_USER" \
  --crdb-database "$CRDB_DATABASE" \
  --redis-host "$REDIS_HOST" \
  --jwt-secret "$JWT_SECRET" \
  --jwt-issuer "$JWT_ISSUER" \
  --jwt-ttl "$JWT_TTL"
