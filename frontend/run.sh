#!/bin/sh

set -eoux pipefail

echo "Dummping envs into file for nodejs"

env > /app/.env.${NODE_ENV:-local}

exec node server.js
