#!/usr/bin/env bash
if [[ "$($PIK simple result)" != *"banana"* ]] ; then
  echo "expected banana" 1>&2
  exit 1
fi