#!/bin/sh

echo "==> Running migrations.."
migrate \
  -source "file://$MIGRATIONS_DIR"  \
  -database "cockroach://$CRDB_USER:$CRDB_PASSWORD@$CRDB_HOST/$CRDB_DATABASE?sslmode=disable" \
  up
