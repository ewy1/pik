#!/usr/bin/env bash
if [[ "$($PIK result)" != *"banana"* ]] ; then
  echo "expected banana" 1>&2
  exit 1
fi