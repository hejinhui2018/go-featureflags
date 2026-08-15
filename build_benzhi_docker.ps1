param(
    [ValidateSet('amd64', 'arm64')]
    [string]$Arch = 'amd64'
)

$ErrorActionPreference = 'Stop'
$tag = "coordinator-p2-limitstore-capacity-update-007:$Arch"
$base = "coordinator-p2-limitstore-capacity-update-007-base:$Arch"

docker build --pull=false --build-arg "BASE_IMAGE=$base" --platform "linux/$Arch" -f benzhi.Dockerfile -t $tag .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
