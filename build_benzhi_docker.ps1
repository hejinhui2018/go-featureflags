param(
    [ValidateSet('amd64', 'arm64')]
    [string]$Arch = 'amd64'
)

$ErrorActionPreference = 'Stop'
$tag = "coordinator-p2-kvjournal-atomic-write-005:$Arch"
$base = "coordinator-p2-kvjournal-atomic-write-005-base:$Arch"

docker build --pull=false --build-arg "BASE_IMAGE=$base" --platform "linux/$Arch" -f benzhi.Dockerfile -t $tag .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
