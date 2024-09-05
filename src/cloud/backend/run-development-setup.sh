#!/bin/bash

set -e

rm -rf data
go build
PROFILE=TEST ./backend
