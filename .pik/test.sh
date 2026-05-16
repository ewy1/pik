#!/usr/bin/env bash
set -euo pipefail
go test -tags test -v ./...
bash .pik/integrations.sh