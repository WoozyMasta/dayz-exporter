#!/usr/bin/env bash
# require upx
set -eu

: "${WORK_DIR:=${1:-cli}}"
: "${BIN_NAME:=dayz-exporter}"

build() {
  local GOOS="${1:-linux}" GOARCH="${2:-amd64}" bin

  bin=$BIN_NAME-$GOOS-$GOARCH
  [ "$GOOS" = windows ] && bin+=.exe

  echo "Build $bin"

  CGO_ENABLED=0 GOARCH="$GOARCH" GOOS="$GOOS" \
  GOFLAGS="-buildvcs=false -trimpath" \
    go build -ldflags="-s -w -X '$version' -X '$commit' -X '$date'" \
      -o "./build/$bin" "$WORK_DIR"/*.go

  [ "$GOOS" = "windows" ] && GOARCH="$GOARCH" go-winres patch --no-backup "./build/$bin"
  [ "$GOOS" = "darwin" ] || [ "$GOOS" = "windows" ] && return

  if command -v xz &>/dev/null; then
    upx --lzma --best "./build/$bin"
    upx -t "./build/$bin"
  fi
}

package="$(grep -Po 'module \K.*$' go.mod)/pkg/config"
version="$package.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo 0.0.0)"
commit="$package.Commit=$(git rev-parse HEAD 2>/dev/null || :)"
date="$package.BuildTime=$(date -uIs)"

mkdir -p ./build
go mod tidy

build darwin amd64
build darwin arm64
build linux 386
build linux amd64
build linux arm
build linux arm64
build windows 386
build windows amd64
build windows arm64
