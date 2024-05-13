#!/bin/bash

set -e

docker pull ocelotcloud/ocelotcloud:demo
docker tag ocelotcloud/ocelotcloud:demo ocelotcloud/ocelotcloud:local
. run.sh