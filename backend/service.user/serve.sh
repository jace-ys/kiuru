#!/bin/sh

exec ./main \
  --port "${PORT:-8080}" \
  --gateway-port "${GATEWAY_PORT:-8081}" \
  --crdb-host "$CRDB_HOST" \
  --crdb-port "${CRDB_PORT:-26257}" \
  --crdb-user "$CRDB_USER" \
  --crdb-dbname "$CRDB_NAME"
