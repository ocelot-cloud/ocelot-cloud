#!/bin/bash

set -e

. .env.sh

docker network create ocelot-net || true

if command -v docker compose &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker-compose"
else
    echo "Error: Neither 'docker compose' nor 'docker-compose' is installed or available in the PATH." >&2
    exit 1
fi

# TODO Just put the docker-compose.yml in this folder
$DOCKER_COMPOSE_CMD -p ocelot-cloud -f "$BACKEND_DIR"/stacks/core/ocelot-cloud/docker-compose.yml up -d

