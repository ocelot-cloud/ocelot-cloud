#!/bin/bash

set -e

go build
./backend -profile="development-setup" -disable-security
