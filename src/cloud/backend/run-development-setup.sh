#!/bin/bash

set -e

rm -rf data
go build
PROFILE=TEST LOG_LEVEL=DEBUG ./backend
