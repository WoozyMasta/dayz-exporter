#!/usr/bin/env bash
set -e

# shellcheck source=/dev/null
. ./functions.sh

geodb::get_country pkg/bemetrics
