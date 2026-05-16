#!/usr/bin/env bash
result="$($PIK defaults dir)"
if [[ ! "$result" == *"apple"* ]] ; then
  echo "expected apple" >&2
  exit 1
fi
