#!/bin/bash

rm -rf data sqlite.db
go build

./hub &
PID=$!
go test -run TestHubRestApi

kill $PID
ls data/users
rm -rf data sqlite.db
