#!/bin/sh

exec ./main \
  --port "$PORT" \
  --gateway-port "$GATEWAY_PORT" \
  --crdb-host "$CRDB_HOST" \
  --crdb-port "${CRDB_PORT:-26257}" \
  --crdb-user "$CRDB_USER" \
  --crdb-dbname "$CRDB_NAME" \
  --crdb-retry "${CRDB_RETRY:-10}"
