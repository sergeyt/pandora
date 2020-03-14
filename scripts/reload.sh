#!/bin/bash

function monitor() {
  if [ "$2" = "true" ];  then
    # Watch all files in the specified directory
    # Call the restart function when they are saved
    inotifywait -q -m -r -e close_write -e moved_to $1 |
    while read line; do
      restart
    done
  else
    # TODO support other file types
    # Watch all *.py files in the specified directory
    # Call the restart function when they are saved
    inotifywait -q -m -r -e close_write -e moved_to --exclude '[^p][^y]$' $1 |
    while read line; do
      restart
    done
  fi
}

# Terminate and rerun the main Go program
function restart {
  if [ "$(pidof $PROCESS_NAME)" ]; then
    killall -q -w -9 $PROCESS_NAME
  fi
  echo ">> Reloading..."
  eval "$ARGS &"
}

# Make sure all background processes get terminated
function close {
  killall -q -w -9 inotifywait
  exit 0
}

trap close INT
echo "== reload"

WATCH_ALL=false
while getopts ":a" opt; do
  case $opt in
    a)
      WATCH_ALL=true
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 0
      ;;
  esac
done

shift "$((OPTIND - 1))"

FILE_NAME=$(basename $1)
PROCESS_NAME=${FILE_NAME%%.*}

ARGS=$@

# Start the main program
echo ">> Watching directories, CTRL+C to stop"
eval "$ARGS &"

monitor $PWD $WATCH_ALL

wait
