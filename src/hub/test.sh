#!/bin/bash

go build
./hub &
PID=$!
go test -run TestCreateUser
kill $PID

rm -rf data