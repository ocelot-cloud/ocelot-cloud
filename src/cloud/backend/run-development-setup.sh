#!/bin/bash

set -e

go build
PROFILE=TEST ./backend -profile="development-setup" -disable-security
