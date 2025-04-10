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

echo "1. Building frontend"
docker build -f Dockerfile.frontend -t frontend-builder:local .
docker run --name frontend-builder -v frontend-deps:/app/node_modules frontend-builder:local
docker cp frontend-builder:/app/dist "$FRONTEND_ARTIFACTS_DIR"
docker rm -f frontend-builder

echo "2. Building backend"
docker build -f Dockerfile.backend -t backend-builder:local .
docker run --name backend-builder -v backend-deps:/go/pkg/mod backend-builder:local .
docker cp backend-builder:/go/src/app/backend "$BACKEND_ARTIFACTS_DIR"
docker rm -f backend-builder

echo "Building production image"
docker build -t ocelotcloud/ocelotcloud:local -f Dockerfile.production .
