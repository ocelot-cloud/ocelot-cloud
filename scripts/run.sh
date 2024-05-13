#!/bin/bash

set -e

docker network create ocelot-net || true
if command docker compose &> /dev/null; then
    docker compose -p ocelot-cloud -f ../components/backend/stacks/core/ocelot-cloud/docker-compose.yml up -d
elif command docker-compose &> /dev/null; then
    docker-compose -p ocelot-cloud -f ../components/backend/stacks/core/ocelot-cloud/docker-compose.yml up -d
else
    echo "Error: Neither 'docker compose' nor 'docker-compose' is installed or available in the PATH." >&2
    exit 1
fi
