# Moodwave CLI — Installation Guide

Moodwave is a native binary. No Node.js, Python, or any runtime required.

---

## Quick Install

### macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex
```

---

## Manual Install

Download the binary for your platform from [GitHub Releases](https://github.com/Boredooms/Moodwave-CLI/releases).

| Platform | File |
|----------|------|
| macOS (Apple Silicon) | `moodwave-darwin-arm64` |
| macOS (Intel) | `moodwave-darwin-amd64` |
| Linux (x86_64) | `moodwave-linux-amd64` |
| Linux (ARM64) | `moodwave-linux-arm64` |
| Linux (ARMv7) | `moodwave-linux-arm` |
| Windows (x86_64) | `moodwave-windows-amd64.exe` |
| Windows (ARM64) | `moodwave-windows-arm64.exe` |

### macOS / Linux

```sh
# Download (replace VERSION and PLATFORM as needed)
curl -fsSL https://github.com/Boredooms/Moodwave-CLI/releases/download/VERSION/moodwave-darwin-arm64 -o moodwave

# Make executable
chmod +x moodwave

# Move to PATH
sudo mv moodwave /usr/local/bin/

# Verify
moodwave --version
```

### Windows

```powershell
# Download the .exe from GitHub Releases
# Move it to a directory in your PATH, for example:
$dest = "$env:LOCALAPPDATA\moodwave\bin"
New-Item -ItemType Directory -Force -Path $dest
Move-Item .\moodwave-windows-amd64.exe "$dest\moodwave.exe"

# Add to PATH (user-level)
$path = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$path;$dest", "User")

# Open new terminal, then verify
moodwave --version
```

---

## Build from Source

Requirements: Go 1.24+

```sh
git clone https://github.com/Boredooms/Moodwave-CLI
cd Moodwave-CLI/cli
go mod download

# Build for current platform
go build -o moodwave ./cmd/moodwave

# Or use Makefile (Unix)
make install

# Or PowerShell (Windows)
scripts\build.ps1 -Only windows
```

---

## Audio Playback Setup

Moodwave streams audio via a subprocess backend. **This step is optional** — if no
backend is found on your `PATH`, Moodwave auto-downloads a static `ffplay` build
into its own cache directory on first playback, and `moodwave doctor` offers a
one-key auto-fix (`[F]`) as well.

To install one yourself instead:

### mpv (recommended — best format support)

```sh
# macOS
brew install mpv

# Ubuntu/Debian
sudo apt install mpv

# Windows
winget install mpv
# or: choco install mpv
```

### ffplay (part of ffmpeg)

```sh
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
winget install ffmpeg
# or: choco install ffmpeg
```

> **Note:** On macOS, `afplay` (CoreAudio) is used as a native fallback. On Windows,
> a native Media Foundation backend is available out of the box. If every backend
> fails, Moodwave degrades to info-only mode — showing the recommended
> station/track and its stream URL so you can open it elsewhere.

---

## Updating

Moodwave updates itself in place:

```sh
moodwave update
```

The interactive TUI also checks for new releases automatically on every launch.
When one is found, a banner appears on the home screen and pressing `[U]`
downloads and installs it without leaving the app.

---

## Verify Installation

Run the built-in diagnostics:

```sh
moodwave doctor
```

This checks:
- Terminal capabilities (color, unicode, TTY)
- Configuration file location
- Cache directory
- Available audio backends
- Music source connectivity (Radio Browser, MusicBrainz, etc.)
- Current session state

---

## Configuration Locations

| Platform | Config file |
|----------|-------------|
| macOS/Linux | `~/.config/moodwave/config.json` |
| Windows | `%APPDATA%\moodwave\config.json` |

| Platform | Cache directory |
|----------|-----------------|
| macOS | `~/Library/Caches/moodwave/` |
| Linux | `~/.cache/moodwave/` |
| Windows | `%LOCALAPPDATA%\moodwave\cache\` |

---

## Uninstall

```sh
# Remove binary
rm $(which moodwave)

# Remove config and cache
rm -rf ~/.config/moodwave ~/.cache/moodwave

# macOS cache
rm -rf ~/Library/Caches/moodwave
```

Windows:
```powershell
Remove-Item "$env:LOCALAPPDATA\moodwave" -Recurse -Force
Remove-Item "$env:APPDATA\moodwave" -Recurse -Force
```
