#!/usr/bin/env bash

set -euo pipefail

# build pik and set $PIK to the built version
go build -o pik .
PIK="$(realpath ./pik)"
export PIK

cd integration_tests

FAILED=""

for dir in *
do
  if [ -d "$dir" ] ; then
    cd "$dir"
    for file in *.test.sh
    do
      tmpdir=$(mktemp -d)
      XDG_CACHE_HOME="$tmpdir"
      XDG_CONFIG_HOME="$tmpdir"
      export XDG_CACHE_HOME XDG_CONFIG_HOME

      if ! bash "$file" 1>/dev/null 2>&1 ; then
        echo "$dir/$file $(tput setaf 1)failed$(tput sgr0)" 2>&1
        bash -x "$file" || true
        FAILED=yes
      else
        echo "$dir/$file $(tput setaf 2)succeeded$(tput sgr0)" 2>&1
        rm -rf "$tmpdir"
      fi
    done
    cd - > /dev/null
  fi
done

if [ -n "$FAILED" ] ; then
  exit 1
fi