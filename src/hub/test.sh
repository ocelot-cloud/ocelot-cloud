#!/bin/bash

go build
./hub &
PID=$!
go test -run TestCreateUser
kill $PID
ls data/users

rm -rf data
