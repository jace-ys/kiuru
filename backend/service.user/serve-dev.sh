#!/bin/sh

exec ./main \
  --port 8080 \
  --gateway-port 8081 \
  --crdb-host localhost \
  --crdb-port 26257 \
  --crdb-user kru_service \
  --crdb-dbname kru \
  --crdb-insecure
