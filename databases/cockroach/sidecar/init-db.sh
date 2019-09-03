#!/bin/sh

HOST=$1

echo "> Connecting to cluster"
until /cockroach/cockroach init --insecure --host=$HOST 2>&1 >/dev/null | grep -q "cluster has already been initialized"; do sleep 1; done

echo "> Initialising cluster"
/cockroach/cockroach sql --insecure --host=$HOST < /cockroach/sidecar/init.sql

echo "> Seeding cluster"
/cockroach/cockroach sql --insecure --host=$HOST < /cockroach/sidecar/seed.sql
