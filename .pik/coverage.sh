#!/usr/bin/env bash
# coverage report
go test -tags test -coverprofile=coverage.out ./...
go tool cover -html coverage.out -o web/coverage.html