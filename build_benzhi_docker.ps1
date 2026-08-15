param(
    [ValidateSet('amd64', 'arm64')]
    [string]$Arch = 'amd64'
)

$ErrorActionPreference = 'Stop'
$tag = "coordinator-p2-cursorstore-validation-008:$Arch"
$binaryDir = Join-Path $PSScriptRoot '.docker-bin'
$binary = Join-Path $binaryDir "cursorstore-$Arch.test"

New-Item -ItemType Directory -Path $binaryDir -Force | Out-Null
$env:GOOS = 'linux'
$env:GOARCH = $Arch
$env:CGO_ENABLED = '0'
$env:GOTOOLCHAIN = 'local'
go test -c -o $binary ./cursorstore
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

docker build --pull=false --build-arg "TARGETARCH=$Arch" --platform "linux/$Arch" -f benzhi.Dockerfile -t $tag .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
