#!/bin/bash

set -e

go build
PROFILE=TEST LOG_LEVEL=DEBUG ./backend
