#!/usr/bin/env bash
set -eu

# Generate test bans.txt for server
# ./tools/random_bans.txt.sh 42 > /path/to/server/bans.txt

: "${COUNT_GUID=${1:-30}}"
: "${COUNT_IP=${2:-20}}"

cd "${0%/*}"
# shellcheck source=/dev/null
. ./functions.sh


for i in $(seq 0 "$COUNT_GUID"); do
  if random::false; then
    time=$(date +%s -ud "$(random::number 1 10080) minute")
  else
    time=-1
  fi
  echo "$(guid::battleye "$(random::steam64)") $time test reason $i"
done

for i in $(seq 0 "$COUNT_IP"); do
  if random::false; then
    time=$(date +%s -ud "$(random::number 1 10080) minute")
  else
    time=-1
  fi
  echo "$(random::ip) $time"
done
