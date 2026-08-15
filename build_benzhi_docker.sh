#!/usr/bin/env sh
set -eu

docker build --platform linux/amd64 -f benzhi.Dockerfile -t coordinator-p2-recordstore-validation-003:amd64 .
docker build --platform linux/arm64 -f benzhi.Dockerfile -t coordinator-p2-recordstore-validation-003:arm64 .
docker run --rm --name coordinator-p2-recordstore-validation-003-amd64 coordinator-p2-recordstore-validation-003:amd64 go test ./...
docker run --rm --name coordinator-p2-recordstore-validation-003-arm64 coordinator-p2-recordstore-validation-003:arm64 go test ./...
