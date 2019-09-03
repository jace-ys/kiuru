#!/bin/sh

HOST=$1

echo "> Initialising cluster"
./cockroach sql --insecure --host=$HOST < /cockroach/sidecar/init.sql

echo "> Seeding cluster"
./cockroach sql --insecure --host=$HOST < /cockroach/sidecar/seed.sql
