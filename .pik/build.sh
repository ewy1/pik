#!/usr/bin/env bash
# build pik
CGO_ENABLED=0 go build -v "$@" .
echo "Congratulations! You just built $(./pik --version)"
