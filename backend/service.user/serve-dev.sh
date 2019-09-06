#!/bin/sh

exec ./main \
  --port "${PORT:-8080}" \
  --gateway-port "${GATEWAY_PORT:-8081}" \
  --crdb-host "${CRDB_HOST:-localhost}" \
  --crdb-port "${CRDB_PORT:-26257}" \
  --crdb-user "${CRDB_USER:-kru_service}" \
  --crdb-dbname "${CRDB_NAME:-kru}" \
  --crdb-retry "${CRDB_RETRY:-10}" \
  --crdb-insecure
