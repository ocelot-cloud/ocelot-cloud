#!/bin/bash

set -e

go build
PROFILE=TEST ./backend
