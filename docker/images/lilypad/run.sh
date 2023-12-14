#!/bin/sh

set -eoux pipefail

if [ -f /app/lilypad/.env ]; then
  source /app/lilypad/.env
fi

exec /usr/local/bin/lilypad $@
