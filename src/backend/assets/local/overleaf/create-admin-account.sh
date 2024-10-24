#!/bin/bash

docker exec overleaf-web /bin/bash -c "grunt user:create-admin --email=admin@admin.com"