# Container verification

Build this image separately for both requested Linux architectures:

```text
docker tag coordinator-p2-featureflags-cancel:amd64 coordinator-p2-recordstore-validation-003-base:amd64
docker tag golang:1.22 coordinator-p2-recordstore-validation-003-base:arm64
docker build --pull=false --build-arg BASE_IMAGE=coordinator-p2-recordstore-validation-003-base:amd64 --platform linux/amd64 -f benzhi.Dockerfile -t coordinator-p2-recordstore-validation-003:amd64 .
docker build --pull=false --build-arg BASE_IMAGE=coordinator-p2-recordstore-validation-003-base:arm64 --platform linux/arm64 -f benzhi.Dockerfile -t coordinator-p2-recordstore-validation-003:arm64 .
docker run --rm --name coordinator-p2-recordstore-validation-003-amd64 coordinator-p2-recordstore-validation-003:amd64 go test ./...
docker run --rm --name coordinator-p2-recordstore-validation-003-arm64 coordinator-p2-recordstore-validation-003:arm64 go test ./...
```
