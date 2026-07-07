# Moodwave-CLI 🌊

> **Your Terminal Mood Music Companion**

Moodwave is an intelligent, developer-focused command-line music player. It automatically scans your local workspace, detects your current working mood based on codebase signals (like TODO density, file count, and language profile), and streams matching music directly from YouTube, Radio Browser, and Jamendo to keep you in the zone.

## Repository Structure

This repository is organized into two main parts:
- **[`cli/`](./cli)**: The core Go implementation of the Moodwave CLI, including source adapters, codebase scanner, terminal visualizers (equalizer/waveforms), and self-updater.
- **[`website/`](./website)**: A beautiful, modern product landing page for Moodwave featuring quick installation commands, screenshots, and guides.

---

## 🚀 Quick Install (No Dependencies Required)

### macOS & Linux
To install Moodwave instantly on macOS or Linux, run:
```bash
curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash
```

### Windows (PowerShell)
To install Moodwave instantly on Windows, run the following in PowerShell:
```powershell
irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex
```

*Note: The installer automatically downloads the pre-compiled binary for your operating system and CPU architecture, installs it, and adds it to your system PATH.*

---

## 🛠️ CLI Development

To build and run the CLI from source:
1. Ensure you have **Go 1.20+** installed.
2. Navigate to the `cli/` folder:
   ```bash
   cd cli
   ```
3. Run the tests:
   ```bash
   go test ./...
   ```
4. Build the binary:
   ```bash
   go build -o moodwave ./cmd/moodwave
   ```

For detailed CLI usage, check out the [CLI Documentation](./cli/docs/commands.md).

## 📄 License
This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.
