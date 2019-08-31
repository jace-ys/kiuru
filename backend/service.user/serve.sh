#!/bin/sh

exec ./main \
  --port "${PORT:-8080}"
  --gateway-port "${GATEWAY_PORT:-8081}"
