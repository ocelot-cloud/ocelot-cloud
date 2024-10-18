#!/bin/bash

rm -rf data
go build
PROFILE="TEST" ./hub
