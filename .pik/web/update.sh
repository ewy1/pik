#!/usr/bin/env bash
# update files on website
set -euo pipefail
ssh git@ewy.one -- cd /srv/pik/pik "&&" bash .pik/web/web.sh