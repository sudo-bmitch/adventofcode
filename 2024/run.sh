#!/bin/sh

if [ $# -lt 1 ]; then
  echo "usage: %0 01a" >&2
  echo "  where 01a is the day and challenge" >&2
  exit 1
fi
day=${1}
go run . ${day} <day${day%?}.input.txt
