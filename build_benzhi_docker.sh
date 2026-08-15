#!/bin/sh
set -eu

platform=${1:?usage: build_benzhi_docker.sh linux/amd64|linux/arm64}
case "$platform" in
  linux/amd64) tag=coordinator-p3-intervalmerge:amd64 ;;
  linux/arm64) tag=coordinator-p3-intervalmerge:arm64 ;;
  *) echo "unsupported platform: $platform" >&2; exit 2 ;;
esac

docker build --platform "$platform" -f benzhi.Dockerfile -t "$tag" .
