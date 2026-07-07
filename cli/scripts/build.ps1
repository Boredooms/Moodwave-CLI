# Moodwave CLI — PowerShell cross-build script
# Runs on Windows without make/bash. Produces binaries for all platforms.
#
# Usage:
#   scripts\build.ps1
#   scripts\build.ps1 -Version v1.0.0
#   scripts\build.ps1 -Only "windows"
param(
    [string]$Version = "dev",
    [string]$Only = "",
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$Module = "github.com/moodwave/moodwave"
$CMD    = "./cmd/moodwave"
$Dist   = "dist"
$BuildTime = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
$LDFlags = "-s -w -X ${Module}/internal/config.Version=$Version -X ${Module}/internal/config.BuildTime=$BuildTime"

# Ensure dist directory.
if (-not (Test-Path $Dist)) {
    New-Item -ItemType Directory -Path $Dist | Out-Null
}

$Targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64";  Ext = ".exe"; Name = "moodwave-windows-amd64.exe" }
    @{ GOOS = "windows"; GOARCH = "arm64";  Ext = ".exe"; Name = "moodwave-windows-arm64.exe" }
    @{ GOOS = "darwin";  GOARCH = "amd64";  Ext = "";     Name = "moodwave-darwin-amd64"      }
    @{ GOOS = "darwin";  GOARCH = "arm64";  Ext = "";     Name = "moodwave-darwin-arm64"       }
    @{ GOOS = "linux";   GOARCH = "amd64";  Ext = "";     Name = "moodwave-linux-amd64"        }
    @{ GOOS = "linux";   GOARCH = "arm64";  Ext = "";     Name = "moodwave-linux-arm64"         }
)

if ($Only) {
    $Targets = $Targets | Where-Object { $_.GOOS -eq $Only }
}

Write-Host ""
Write-Host "  Moodwave CLI — Cross Compile" -ForegroundColor White
Write-Host ("  " + "─" * 40)
Write-Host "  Version:    $Version"
Write-Host "  Build time: $BuildTime"
Write-Host "  Output:     $Dist/"
Write-Host ""

$Success = 0
$Failed = 0

foreach ($t in $Targets) {
    $OutFile = Join-Path $Dist $t.Name
    $env:GOOS = $t.GOOS
    $env:GOARCH = $t.GOARCH
    $env:CGO_ENABLED = "0"

    Write-Host "  Building $($t.Name)..." -NoNewline

    try {
        if ($Verbose) {
            & go build -v -ldflags $LDFlags -o $OutFile $CMD
        }
        else {
            & go build -ldflags $LDFlags -o $OutFile $CMD 2>&1 | Out-Null
        }
        if ($LASTEXITCODE -eq 0) {
            $Size = (Get-Item $OutFile).Length / 1024
            Write-Host " OK ($([math]::Round($Size))KB)" -ForegroundColor Green
            $Success++
        }
        else {
            Write-Host " FAILED" -ForegroundColor Red
            $Failed++
        }
    }
    catch {
        Write-Host " ERROR: $_" -ForegroundColor Red
        $Failed++
    }
}

# Cleanup env.
Remove-Item env:GOOS -ErrorAction SilentlyContinue
Remove-Item env:GOARCH -ErrorAction SilentlyContinue
Remove-Item env:CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "  Built: $Success / $($Success + $Failed) targets" -ForegroundColor $(if ($Failed -eq 0) { "Green" } else { "Yellow" })

if ($Success -gt 0) {
    Write-Host ""
    Write-Host "  Output files:"
    Get-ChildItem $Dist | Select-Object Name, @{N='Size(KB)';E={[math]::Round($_.Length/1024)}} | Format-Table -AutoSize
}

if ($Failed -gt 0) {
    exit 1
}
