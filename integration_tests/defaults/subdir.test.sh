#!/usr/bin/env bash
result="$($PIK dir)"
if [[ ! "$result" == *"apple"* ]] ; then
  echo "expected apple" >&2
  exit 1
fi
