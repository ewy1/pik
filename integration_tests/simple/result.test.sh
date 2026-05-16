#!/usr/bin/env bash
if [[ "$($PIK echo asdf)" != *"asdf"* ]] ; then
  echo "expected banana" 1>&2
  exit 1
fi