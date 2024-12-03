#!/usr/bin/env bash
set -eu

# Generate test data for RCON responce of players command
# ./tools/random_rcon_players.sh 42 > ./pkg/bemetrics/test_data/players.txt

: "${COUNT=${1:-50}}"

cd "${0%/*}"
# shellcheck source=/dev/null
. ./functions.sh

echo 'Players on server:'
echo '[#] [IP Address]:[Port] [Ping] [GUID] [Name]'
echo '--------------------------------------------------'

i=0
while read -r name; do
  guid="$(guid::battleye "$(random::steam64)")"
  name="${name//[$'\t\r\n']}$(random::lobby)"

  printf "%-2s %-23s %-5s %-12s(OK) %-16s\n" \
    "$i" "$(random::ip_port)" "$(random::number 1 300)" "${guid^^}" "$name"

  (( i++ )) || :
done < <(random::users "$COUNT")

echo "($COUNT players in total)"
