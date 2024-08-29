#!/bin/bash

set -e

. .env.sh

docker network create ocelot-net || true
$DOCKER_COMPOSE_CMD -p ocelot-cloud -f "$BACKEND_DIR"/stacks/ocelot-cloud/docker-compose.yml up -d

