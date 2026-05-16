#!/usr/bin/env bash
result="$($PIK defaults defaults)"
if [[ ! "$result" == *"banana"* ]] ; then
  echo "expected banana" >&2
  exit 1
fi