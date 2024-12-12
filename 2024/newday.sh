#!/bin/sh

num="${1}"
if ! ( echo $num | grep '^[0-9][0-9]$' >/dev/null ); then 
  echo "invalid day, must be a two digit number: $num"
  exit 1
fi

if [ -f "day${num}.go" ] || [ -f "day${num}_test.go" ]; then
  echo "day already created"
  exit 1
fi

if ! ( [ -f "dayXX.go" ] && [ -f "dayXX_test.go" ] ); then
  echo "template days are missing"
  exit 1
fi

sed "s/XX/${num}/g" <dayXX.go      >day${num}.go
sed "s/XX/${num}/g" <dayXX_test.go >day${num}_test.go
