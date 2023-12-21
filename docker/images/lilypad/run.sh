#!/bin/sh

set -eoux pipefail

# Create bacalhau config
/usr/local/bin/bacalhau version

if [ -f /app/lilypad/.env ]; then
  source /app/lilypad/.env
fi

exec /usr/local/bin/lilypad $@
