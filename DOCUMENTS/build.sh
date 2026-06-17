#!/bin/bash
set -e

NO_CACHE=""
[ "${1:-}" = "--no-cache" ] && NO_CACHE="--no-cache"

docker build $NO_CACHE \
    -f docker/nagios-core/nagios/Dockerfile \
    -t nagios-core:latest \
    docker/nagios-core/
