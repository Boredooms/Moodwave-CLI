#Requires -Version 5.1
<#
.SYNOPSIS
    Moodwave CLI installer for Windows.
.DESCRIPTION
    Downloads and installs the Moodwave CLI binary from GitHub Releases.
    Adds the install directory to the user's PATH.
.EXAMPLE
    irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex
.EXAMPLE
    $env:MOODWAVE_VERSION = "v1.0.0"; irm https://... | iex
#>

param(
    [string]$Version = $env:MOODWAVE_VERSION,
    [string]$InstallDir = "$env:LOCALAPPDATA\moodwave\bin"
)

$ErrorActionPreference = "Stop"
$Repo = "Boredooms/Moodwave-CLI"
$Binary = "moodwave"

function Write-Info    { param($msg) Write-Host "  info  " -ForegroundColor Cyan -NoNewline; Write-Host $msg }
function Write-Success { param($msg) Write-Host "  ok    " -ForegroundColor Green -NoNewline; Write-Host $msg }
function Write-Warn    { param($msg) Write-Host "  warn  " -ForegroundColor Yellow -NoNewline; Write-Host $msg }
function Write-Fail    { param($msg) Write-Host "  error " -ForegroundColor Red -NoNewline; Write-Host $msg; exit 1 }

# ── Banner ────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  Moodwave CLI Installer" -ForegroundColor White
Write-Host ("  " + "─" * 40)
Write-Host ""

# ── Detect architecture ───────────────────────────────────────────────────────
$Arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "amd64" }
    "ARM64" { "arm64" }
    default { "amd64" }
}

# ── Resolve version ───────────────────────────────────────────────────────────
if (-not $Version -or $Version -eq "latest") {
    Write-Info "Resolving latest version..."
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        $Version = $release.tag_name
    }
    catch {
        Write-Fail "Could not resolve latest version: $_. Set `$env:MOODWAVE_VERSION`."
    }
}

Write-Info "Version:    $Version"
Write-Info "Arch:       $Arch"
Write-Info "Install to: $InstallDir"
Write-Host ""

# ── Build download URL ────────────────────────────────────────────────────────
$BinaryFile = "$Binary-windows-$Arch.exe"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$BinaryFile"

# ── Create install directory ──────────────────────────────────────────────────
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# ── Download ──────────────────────────────────────────────────────────────────
$TmpFile = Join-Path $env:TEMP "moodwave-download.exe"
Write-Info "Downloading $BinaryFile..."
try {
    $ProgressPreference = "SilentlyContinue"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TmpFile -UseBasicParsing
}
catch {
    Write-Fail "Download failed from: $DownloadUrl`nError: $_"
}

# ── Verify binary ─────────────────────────────────────────────────────────────
try {
    $testOut = & $TmpFile --version 2>&1
    Write-Info "Verified: $testOut"
}
catch {
    Remove-Item $TmpFile -Force -ErrorAction SilentlyContinue
    Write-Fail "Downloaded binary failed to run: $_"
}

# ── Install ───────────────────────────────────────────────────────────────────
$InstallPath = Join-Path $InstallDir "$Binary.exe"
Copy-Item $TmpFile $InstallPath -Force
Remove-Item $TmpFile -Force

Write-Success "Installed to $InstallPath"

# ── Add to PATH ───────────────────────────────────────────────────────────────
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
    $env:PATH = "$env:PATH;$InstallDir"
    Write-Success "Added $InstallDir to PATH (restart terminal to take effect)"
}
else {
    Write-Info "$InstallDir already in PATH"
}

# ── Verify installation ───────────────────────────────────────────────────────
try {
    $installedVersion = & $InstallPath --version 2>&1
    Write-Success "Installed: $installedVersion"
}
catch {
    Write-Warn "Binary installed but could not be verified in PATH"
}

# ── Final instructions ────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  Ready! Open a new terminal and try:" -ForegroundColor White
Write-Host "    moodwave init"
Write-Host "    moodwave scan"
Write-Host "    moodwave play"
Write-Host ""
