#!/usr/bin/env bash
# install pik to ~/.local/bin
set -euo pipefail
DEST_FOLDER="$HOME/.local/bin"
if [ ! -d "$DEST_FOLDER" ] ; then
  echo "I only know how to install to ~/.local/bin, sorry!"
fi
DEST="$DEST_FOLDER/pik"
if [ -f "$DEST" ] ; then
  rm "$DEST"
fi

go build -o "$DEST" "$@" .
echo "congratulations!"