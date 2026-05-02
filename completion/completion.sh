#/usr/bin/env bash
_pik_completions()
{
  QUERY=""
  for word in COMP_WORDS ; do
    if [ ! query = "-"* ] ; then
      QUERY="$QUERY $WORD"
    fi
  done
  COMPREPLY=($(compgen -W "$(pik --list)" "${QUERY}"))
}

complete -F _pik_completions pik