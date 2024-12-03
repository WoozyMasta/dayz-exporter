#!/usr/bin/env bash
# Some useful functions

(return 0 2>/dev/null) || >&2 printf '\n\t%s\n\t%s\n\n' \
  'This file is not a script that executes logic, it is just a set of functions' \
  "Use it as \"source $0\""

# Generate Random Steam64 ID
# lowest SteamID is 76561197960265728 (0x0110000100000000)
# higest SteamID is 76561202255233023 (0x01100001FFFFFFFF)
# total 4294967295
random::steam64() {
  local min=197960265728 max=202255233023

  rnd=$((RANDOM % (max - min + 1) + min))
  echo "76561$rnd"
}

# Generate BattleEye GUID from Steam64 ID
# md5("BE" (2 bytes) + 64-bit SteamID (8 bytes))
# example: `guid::battleye "$(random::steam64)"`
guid::battleye () {
  printf 'BE%s' "$1" | md5sum | awk '{print toupper($1)}'
}

# Generate DayZ GUID from Steam64 ID
# Base64(SHA256(SteamID64))
# example: `guid::dayz "$(random::steam64)"`
guid::dayz () {
  printf '%s' "$1" | sha256sum | awk '{print $1}' | base64 -w0
}

# Generate random IP:PORT string
random::ip_port() {
  printf '%d:%d' "$(random::ip)" "$(( ( RANDOM % 65535 ) + 20000 ))"
}

# Generate random IP
random::ip() {
  printf '%d.%d.%d.%d' \
    "$(( ( RANDOM % 256 ) + 1 ))" "$(( ( RANDOM % 256 ) + 1 ))" \
    "$(( ( RANDOM % 256 ) + 1 ))" "$(( ( RANDOM % 256 ) + 1 ))"
}

# Generate random number in range
random::number() {
  local min=${1:-1} max=${2:-300}
  printf '%d' "$(( ( RANDOM % max ) + min ))"
}

# Random flip 1/10 string (Lobby)
random::lobby() {
  local chance=${1:-10}
  if [ $(( RANDOM % chance )) -eq 5 ]; then
    printf ' (Lobby)'
  fi
}

# Random flip, return false (1) in chance
random::false() {
  local chance=${1:-10}
  if [ $(( RANDOM % chance )) -eq 5 ]; then
    return 1
  fi
}

# Get random user names from randomuser.me API
random::users() {
  local count=${1:-50}
  curl -sSfL \
    "https://randomuser.me/api/?results=$count&nat=ua&gender=male&inc=name" | \
  jq -er '.results[].name | [.first, .last] | join(" ")' | \
  iconv -c -f utf-8 -t ascii
}

# Download geo DB
geodb::get() {
  local path="${1:-.}" db="${2:-GeoLite2-Country.mmdb}"
  echo "Download $db to $path/"
  curl -#SfLo "$path/$db" "https://git.io/$db"
}

# Download GeoLite2-City.mmdb
geodb::get_city() {
  geodb::get "${1:-.}" GeoLite2-City.mmdb
}

# Download GeoLite2-Country.mmdb
geodb::get_country() {
  geodb::get "${1:-.}" GeoLite2-Country.mmdb
}

# Download GeoLite2-ASN.mmdb
geodb::get_asn() {
  geodb::get "${1:-.}" GeoLite2-ASN.mmdb
}
