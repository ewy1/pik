#!/usr/bin/env bash
# build uwu
echo "$@"
CGO_ENABLED=0 go build -v "$@" .
echo "Congratulations!"
