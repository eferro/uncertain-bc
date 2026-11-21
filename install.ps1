# Installs ubc on Windows:
#   irm https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.ps1 | iex
# $env:UBC_VERSION pins a release (default: latest); $env:UBC_INSTALL_DIR picks
# the directory (default: %LOCALAPPDATA%\Programs\ubc, added to the user PATH).
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$repo = 'eferro/uncertain-bc'
$version = $env:UBC_VERSION
if (-not $version) {
    $version = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
}
$version = $version.TrimStart('v')
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
$dir = if ($env:UBC_INSTALL_DIR) { $env:UBC_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\ubc' }
$archive = "ubc_${version}_windows_$arch.zip"
$base = "https://github.com/$repo/releases/download/v$version"

$work = New-Item -ItemType Directory -Path (Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid()))
try {
    Invoke-WebRequest "$base/$archive" -OutFile (Join-Path $work $archive)
    Invoke-WebRequest "$base/checksums.txt" -OutFile (Join-Path $work 'checksums.txt')
    $line = Get-Content (Join-Path $work 'checksums.txt') | Where-Object { $_ -match " $([regex]::Escape($archive))$" }
    if (-not $line) { throw "no checksum for $archive" }
    $expected = ($line -split '\s+')[0]
    $actual = (Get-FileHash (Join-Path $work $archive) -Algorithm SHA256).Hash.ToLower()
    if ($actual -ne $expected) { throw "checksum mismatch for $archive" }

    Expand-Archive (Join-Path $work $archive) -DestinationPath $work -Force
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    Copy-Item (Join-Path $work 'ubc.exe') (Join-Path $dir 'ubc.exe') -Force
    # A downloaded file carries the Mark of the Web, which makes SmartScreen
    # warn about the unsigned binary; the checksum above already verified it.
    Unblock-File (Join-Path $dir 'ubc.exe')
} finally {
    Remove-Item $work -Recurse -Force
}

Write-Host "ubc $version installed to $(Join-Path $dir 'ubc.exe')"
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($userPath -split ';') -notcontains $dir) {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$dir", 'User')
    Write-Host "Added $dir to your PATH; open a new terminal to use ubc."
}
