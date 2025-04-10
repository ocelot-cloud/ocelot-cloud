#!/bin/sh

# Take some time to get ready and make port available.
sleep 1

terminate() {
  echo "Shutting down Nginx..."
  nginx -s quit
  exit 0
}

trap 'terminate' TERM
nginx -g 'daemon off;' &
wait $!
