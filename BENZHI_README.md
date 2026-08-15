# Container verification

Build one architecture at a time with `./build_benzhi_docker.ps1 -Arch amd64` or `./build_benzhi_docker.ps1 -Arch arm64`. The script cross-compiles the package tests with the local Go toolchain and creates a minimal image for the selected Linux architecture.

The image runs the compiled Go test suite by default and does not expose a host port.
