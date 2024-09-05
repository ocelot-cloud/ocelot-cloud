#!/bin/bash

set -e

rm -rf data
go build
PROFILE=TEST USE_DUMMY_STACKS=true ./backend
