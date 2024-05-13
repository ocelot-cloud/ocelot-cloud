#!/bin/sh

# Take some time to get ready and make port available.
sleep 3

# Close port immediately, but wait a few seconds before the running container can be destroyed by docker.
terminate() {
  echo "Shutting down Nginx..."
  nginx -s quit
  sleep 2
  exit 0
}

trap 'terminate' TERM
nginx -g 'daemon off;' &
wait $!
