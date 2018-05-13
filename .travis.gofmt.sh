#!/bin/bash

result="$(gofmt -s -l . | grep -v '^vendor/' )"
if [ -n "$result" ]; then
  echo "Go code is not formatted, run 'gofmt -s -w .'" >&2
  echo "$result"
  exit 1
fi
