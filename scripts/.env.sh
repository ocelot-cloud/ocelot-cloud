#!/bin/bash

echo "Setting directory variables"
cd ../src
SRC_DIR="$(pwd)"
DOCKER_DIR="$SRC_DIR/ci-runner/docker"
FRONTEND_DIR="$SRC_DIR/cloud/frontend"
BACKEND_DIR="$SRC_DIR/cloud/backend"
ARTIFACTS_DIR="$DOCKER_DIR/artifacts"
FRONTEND_ARTIFACTS_DIR="$ARTIFACTS_DIR/frontend"
BACKEND_ARTIFACTS_DIR="$ARTIFACTS_DIR/backend"