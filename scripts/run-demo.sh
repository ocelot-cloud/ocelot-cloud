#!/bin/bash

set -e

docker pull ocelotcloud/ocelotcloud:demo
docker tag ocelotcloud/ocelotcloud:demo ocelotcloud/ocelotcloud:local
docker network create ocelot-net || true
docker compose -p ocelot-cloud -f ../src/backend/assets/ocelot-cloud/docker-compose.yml up -d