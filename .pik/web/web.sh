#!/usr/bin/env bash
# gets run on server after calling `pik web update` from anywhere
set -euo pipefail
git pull
bash .pik/coverage.sh