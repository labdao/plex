#!/bin/sh
set -e

find /container-init.d -maxdepth 1 -type f -iname '*.sh' -print0 | sort -z | xargs -n 1 -0 -r container_init_run

exec /plex "$@"
