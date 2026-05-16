#!/usr/bin/env bash
if [[ "$($PIK 123)" != "" ]] ; then
  echo "expected no output" 1>&2
  exit 1
fi