# Container verification

Build one architecture at a time with `./build_benzhi_docker.ps1 -Arch amd64` or `./build_benzhi_docker.ps1 -Arch arm64` after preparing the matching private base tag named in the script.

The image runs `go test ./...` by default and does not expose a host port.
