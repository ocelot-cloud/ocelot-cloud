#!/bin/bash

set -e

echo "Cleaning up"
docker rm -f frontend-builder
docker rm -f backend-builder
docker rm -f ocelotcloud/ocelotcloud

echo "Creating necessary docker network and volumes"
docker network create ocelot-net || true
docker volume create backend-deps
docker volume create frontend-deps

. .env.sh

echo "Transferring source code to artifacts directories"
mkdir -p "$FRONTEND_ARTIFACTS_DIR" "$BACKEND_ARTIFACTS_DIR"
rsync -a --delete --exclude="node_modules" "$FRONTEND_DIR" "$ARTIFACTS_DIR/"
rsync -a --delete "$BACKEND_DIR" "$ARTIFACTS_DIR/"

echo "Building components"
cd "$DOCKER_DIR"

# TODO Resolve duplication of the docker images, there are two times here and once in the Dockerfiles each
echo "1. Building frontend"
if [ -z "$(docker images -q node:18.10.0)" ]; then docker pull node:18.10.0; fi
docker build -f Dockerfile.frontend -t frontend-builder:local .
docker run --name frontend-builder -v frontend-deps:/app/node_modules frontend-builder:local
docker cp frontend-builder:/app/dist "$FRONTEND_ARTIFACTS_DIR"
docker rm -f frontend-builder

echo "2. Building backend"
if [ -z "$(docker images -q golang:1.21.8)" ]; then docker pull golang:1.21.8; fi
docker build -f Dockerfile.backend -t backend-builder:local .
docker run --name backend-builder -v backend-deps:/go/pkg/mod backend-builder:local .
docker cp backend-builder:/go/src/app/backend "$BACKEND_ARTIFACTS_DIR"
docker rm -f backend-builder

echo "Building production image"
if [ -z "$(docker images -q alpine:3.18.6)" ]; then docker pull alpine:3.18.6; fi
docker build -t ocelotcloud/ocelotcloud:local -f Dockerfile.production .

# TODO Problem: using the docker builders makes problems when there is no internet connection. For development I should build the stuff natively. Also remove that from the README. Instead tell to use "install.sh" and then "ci-runner build" or so.