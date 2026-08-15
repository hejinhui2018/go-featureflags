param(
    [ValidateSet('amd64', 'arm64')]
    [string]$Arch = 'amd64'
)

$ErrorActionPreference = 'Stop'
$tag = "coordinator-p2-ttlstore-duration-validation-009:$Arch"
$binaryDir = Join-Path $PSScriptRoot '.docker-bin'
$binary = Join-Path $binaryDir "ttlstore-$Arch.test"

New-Item -ItemType Directory -Path $binaryDir -Force | Out-Null
$env:GOOS = 'linux'
$env:GOARCH = $Arch
$env:CGO_ENABLED = '0'
$env:GOTOOLCHAIN = 'local'
go test -c -o $binary ./ttlstore
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

docker build --pull=false --build-arg "TARGETARCH=$Arch" --platform "linux/$Arch" -f benzhi.Dockerfile -t $tag .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
