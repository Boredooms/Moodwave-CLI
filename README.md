# Moodwave-CLI 🌊

> **A terminal-native developer mood music companion.**

<p align="left">
  <img src="./website/public/logo.svg" width="80" height="80" alt="Moodwave Logo" />
</p>

Moodwave scans your codebase, extracts cognitive and workspace signals, infers your current coding mood profile, and streams perfectly matched audio directly inside your shell. Featuring real-time ANSI waveforms, multiple visualizer modes, auto-updates, and native multi-platform compilation—built entirely for developers who don't want to think about what to play.

---

## ⚡ Quick Install (Zero Dependencies)

No Go runtime, Python, or external languages required. Runs natively out of the box.

### macOS & Linux
```bash
curl -fsSL https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Boredooms/Moodwave-CLI/main/cli/scripts/install.ps1 | iex
```

*The installer automatically fetches the optimized binary matching your CPU (amd64, arm64, or arm) and OS, placing it in your system PATH.*

---

## 🏗️ Technical Architecture

Moodwave is designed as a decoupled, multi-layered streaming client:

```
                  ┌─────────────────────────────────────┐
                  │          Codebase Workspace         │
                  └──────────────────┬──────────────────┘
                                     │ (Files, Git, TODOs)
                                     ▼
                  ┌─────────────────────────────────────┐
                  │ 01. Workspace Scanner               │
                  └──────────────────┬──────────────────┘
                                     │ (Language metrics, git tree depth)
                                     ▼
                  ┌─────────────────────────────────────┐
                  │ 02. Weighted Heuristics Mood Engine │
                  └──────────────────┬──────────────────┘
                                     │ (10 Developer Mood Profiles)
                                     ▼
                  ┌─────────────────────────────────────┐
                  │ 03. Tag-Based Recommender           │
                  └──────────────────┬──────────────────┘
                                     │ (BPM range, Tag overlap, History)
                                     ▼
         ┌───────────────────────────┼───────────────────────────┐
         ▼                           ▼                           ▼
┌──────────────────┐       ┌──────────────────┐        ┌──────────────────┐
│ YouTube (yt-dlp) │       │ Radio Browser API│        │   Jamendo API    │
└────────┬─────────┘       └────────┬─────────┘        └────────┬─────────┘
         │                          │                           │
         └──────────────────────────┼───────────────────────────┘
                                    ▼
                  ┌─────────────────────────────────────┐
                  │ 04. Playback Controller (mpv/ffplay)│
                  └──────────────────┬──────────────────┘
                                     │ (HTTP Range Resume & Stream)
                                     ▼
                  ┌─────────────────────────────────────┐
                  │ 05. Visual Renderer (ANSI TUI)      │
                  └──────────────────┬──────────────────┘
                                     │ (Waveform, Equalizer, Spectrum)
                                     ▼
                  ┌─────────────────────────────────────┐
                  │          Developer Terminal         │
                  └─────────────────────────────────────┘
```

1. **Workspace Scanner**: Crawls the workspace to analyze language composition, TODO/FIXME comments density, git tree activity, and file count variance to build a codebase signature profile.
2. **Mood Inference Engine**: Uses a weighted heuristics rule chain to map workspace signals to one of 10 developer moods (e.g., `focused`, `calm`, `intense`, `chaotic`, `sprint`, `debugging`).
3. **Recommender**: Selects and ranks candidate audio tracks matching the inferred mood parameters (BPM bounds, genre tags, energy levels, and history state).
4. **Playback Controller**: Automatically dispatches streaming audio threads using system audio players (`mpv` or `ffplay`) with automated network re-establishment.
5. **ANSI visualizer**: Dynamically renders real-time visualizers (sine waveforms, spectrum bars, pulses) directly inside the shell using ANSI color escape sequences.

---

## 🛠️ CLI Command Reference

| Command | Description |
|---|---|
| `moodwave init` | Initialize default TOML configuration and directories |
| `moodwave scan` | Force-scan the repository and show inferred mood metrics |
| `moodwave play` | Play streams matching your current coding mood |
| `moodwave search` | Query YouTube manually and play the selected stream |
| `moodwave status` | Print current CLI and playback progress information |
| `moodwave next` | Skip the current track and fetch the next recommendation |
| `moodwave theme` | Toggle active TUI display themes dynamically |
| `moodwave visual` | Switch visualizer modes (waveform, spectrum, minimal, quiet) |
| `moodwave update` | Download and replace the binary with the latest release |
| `moodwave doctor` | Inspect and verify local system capabilities and audio components |

---

## 💻 Local Compilation & Development

Ensure you have **Go 1.22+** installed on your development machine.

1. Clone the repository:
   ```bash
   git clone https://github.com/Boredooms/Moodwave-CLI.git
   cd Moodwave-CLI/cli
   ```
2. Run unit tests:
   ```bash
   go test ./...
   ```
3. Compile the local binary:
   ```bash
   go build -o moodwave ./cmd/moodwave
   ```
4. Verify the executable:
   ```bash
   ./moodwave doctor
   ```

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for details.
