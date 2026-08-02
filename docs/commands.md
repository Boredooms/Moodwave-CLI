# Moodwave CLI — Command Reference

Complete reference for all `moodwave` commands.

---

## Global Flags

Available on all commands:

| Flag | Description |
|------|-------------|
| `--help`, `-h` | Show help for this command |
| `--version`, `-v` | Show version and build info |
| `--debug` | Enable verbose debug logging |
| `--no-color` | Disable ANSI color codes |
| `--no-animation` | Disable terminal animations |
| `--path <dir>` | Override the scanned directory (default: current) |

---

## moodwave init

Initialize Moodwave config and create required directories.

```sh
moodwave init
```

Creates:
- Config file at `~/.config/moodwave/config.json` (or platform equivalent)
- Cache directory at `~/.cache/moodwave/` (or platform equivalent)
- Default configuration with safe built-in values

Run once after installing. Safe to run again — will not overwrite existing config.

---

## moodwave scan

Scan the current repository and detect the coding mood.

```sh
moodwave scan
moodwave scan /path/to/repo
moodwave scan --path /path/to/repo
```

**What it does:**
1. Walks the repository directory tree (respects standard ignore dirs)
2. Extracts codebase signals:
   - Language distribution and primary language
   - Test file ratio and coverage
   - TODO/FIXME density
   - Comment density
   - Directory structure entropy
   - Naming consistency
   - Build system and dependencies
   - CI configuration
   - Git branch, churn score, last commit age
3. Runs the mood inference engine (13 heuristic rules)
4. Saves the mood profile to the session cache

**Output:** Mood label, confidence percentage, explanation, BPM range, and music tags.

**Options:**

| Flag | Description |
|------|-------------|
| `--path <dir>` | Directory to scan (default: current working directory) |

---

## moodwave mood

Show the current mood profile without re-scanning.

```sh
moodwave mood
```

Displays:
- Mood label and emoji
- Confidence percentage
- Age of last scan
- Explanation
- Music traits (BPM range, energy, tags)
- Top contributing signals
- All 10 mood scores as a bar chart

Requires a previous `moodwave scan` to have been run.

---

## moodwave play

Find and play music matching the detected mood.

```sh
moodwave play
```

**What it does:**
1. Loads the current mood profile (or runs `scan` if none exists)
2. Queries all healthy music sources in priority order:
   - Radio Browser (stations by mood tags)
   - Jamendo (CC tracks, if client_id configured)
   - MusicBrainz (metadata enrichment)
3. Scores all candidates against mood traits (BPM, energy, tag overlap)
4. Starts the best audio backend (mpv → ffplay → afplay → PowerShell)
5. Launches the terminal visual renderer

**Visual controls during playback (when TTY):**
- See keyboard shortcuts at the bottom of the screen
- Press `Q` or `Ctrl-C` to stop

**No audio backend?** Displays the stream URL for manual playback.

---

## moodwave pause

Signal playback to pause (requires active play session in same terminal).

```sh
moodwave pause
```

> **Note:** Since `moodwave play` manages its own process, pause is most effective
> as `Ctrl-C` during the play session, then restarting with `moodwave play`.
> Full IPC-based pause is planned for a future version.

---

## moodwave stop

Mark the session as stopped.

```sh
moodwave stop
```

---

## moodwave next

Advance to the next candidate in the recommendation queue.

```sh
moodwave next
```

Run `moodwave play` to start playing the next track.

---

## moodwave queue

Show the current recommendation queue.

```sh
moodwave queue
```

Displays each candidate with:
- Position (▶ marks the current)
- Display name
- Score percentage
- Recommendation reason

---

## moodwave status

Show the current CLI and playback status.

```sh
moodwave status
```

Displays:
- Version and build info
- Config and cache locations
- Current project path
- Last detected mood and age
- Playback state
- Current track/station
- Terminal capabilities

---

## moodwave config

View or edit the configuration.

```sh
# Show full config as JSON
moodwave config

# Show config file path
moodwave config path

# Open in editor ($EDITOR or default)
moodwave config edit
```

Config is stored as JSON. Key sections:

```json
{
  "scanner": {
    "max_depth": 12,
    "max_files": 50000,
    "git_enabled": true
  },
  "sources": {
    "priority": ["radio_browser", "jamendo", "musicbrainz"],
    "radio_browser": { "enabled": true },
    "jamendo": { "enabled": false, "client_id": "" }
  },
  "playback": {
    "backend": "",
    "volume": 80
  },
  "visual": {
    "mode": "wave",
    "theme": "monochrome",
    "fps": 24
  }
}
```

---

## moodwave theme

Switch the visual color theme.

```sh
# List available themes
moodwave theme

# Switch theme
moodwave theme dark
moodwave theme ash
moodwave theme ghost
moodwave theme monochrome
```

**Available themes:**

| Theme | Description |
|-------|-------------|
| `midnight` | Deep blue-black, high contrast (TUI default) |
| `dracula` | Purple/pink on dark gray |
| `nord` | Cool arctic blues |
| `tokyo-night` | Muted indigo and cyan |
| `gruvbox` | Warm retro earth tones |
| `catppuccin` | Soft pastel on mocha |
| `solarized` | Balanced low-contrast palette |
| `monochrome` | Pure grayscale |
| `dark` | Dark with subtle blue tints |
| `ash` | Warm gray |
| `ghost` | Near-invisible, minimal contrast |

> In the interactive TUI, press `[T]` to open the theme selector with a live
> color preview. Themes apply immediately across every screen.

---

## moodwave visual

Switch the visual animation mode.

```sh
# List available modes
moodwave visual

# Switch mode
moodwave visual wave
moodwave visual spectrum
moodwave visual pulse
moodwave visual minimal
moodwave visual quiet
```

**Available modes:**

| Mode | Description |
|------|-------------|
| `spectrum` | Equalizer bars (TUI default) |
| `fire` | Doom-style cellular-automaton fire |
| `matrix` | Falling glyph rain |
| `rain` | Vertical rainfall |
| `plasma` | Smooth interference field |
| `aurora` | Drifting polar light bands |
| `snow` | Slow drifting snowfall |
| `fireflies` | Scattered wandering glow points |
| `lava` | Rising molten blobs |
| `dna` | Rotating double helix |
| `campfire` | Warmer, narrower fire |
| `smoke` | Rising diffuse plumes |
| `ripple` | Concentric water rings |
| `lightning` | Intermittent forked strikes |
| `spiral` | Rotating logarithmic spiral |
| `petals` | Falling flower petals |
| `galaxy` | Star field with drift |
| `zen` | Slow calm ocean swell |
| `heartbeat` | ECG-style pulse trace |
| `pendulum` | Swinging pendulum arc |
| `city` | Night skyline silhouette |
| `ocean` | Rolling ocean waves |
| `starfield` | Forward-motion star travel |
| `pulse` | Ambient pulse rings |
| `minimal` | Status text only, no animation |
| `quiet` | No visual panel at all |

> In the interactive TUI, press `[V]` to cycle visualizers live. All modes
> size themselves to the current terminal width and height.

---

## moodwave source

List music sources and their health status.

```sh
# List all sources with health check
moodwave source
moodwave source list
```

**Adding optional sources:**

```sh
# Jamendo (CC music catalog)
export JAMENDO_CLIENT_ID=your_client_id
moodwave play

# ListenBrainz (personalized recommendations)
export LISTENBRAINZ_TOKEN=your_user_token
export LISTENBRAINZ_USERNAME=your_username
moodwave play
```

---

## moodwave doctor

Run a complete system diagnostic.

```sh
moodwave doctor
```

Checks:
1. **System** — OS, architecture
2. **Terminal** — TTY, color, unicode, animation support
3. **Configuration** — config file and cache directory presence
4. **Audio backends** — mpv, ffplay, afplay, vlc availability with install hints
5. **Music sources** — live health check for each adapter
6. **Cache** — entry count and limits
7. **Session** — last scan age and mood

In a TTY this runs as an interactive TUI. Press `[F]` to auto-fix a failing
check — most commonly downloading a static `ffplay` build when no audio backend
is available. In a non-TTY (piped/CI) context it falls back to plain text output.

---

## moodwave update

Update the binary in place to the latest GitHub release.

```sh
moodwave update
# alias: moodwave upgrade
```

**What it does:**
1. Queries the GitHub releases API for the latest tag
2. Compares it against the running binary's version — exits early if already current
3. Picks the release asset matching your OS and architecture
4. Downloads it and atomically swaps the running executable (with rollback if the write fails)

The interactive TUI performs the same check automatically on every launch. When
an update is available, a banner appears on the home screen and pressing `[U]`
installs it without leaving the app. A failed background check is silent — it
never interrupts you with an error.

---

## Interactive TUI

Running `moodwave` with no arguments in a TTY launches the full interactive
interface. Screens: Home, Search, Now Playing, Live Queue, Personal Playlist,
Theme Select, and Doctor.

### Keyboard shortcuts

| Key | Context | Action |
|---|---|---|
| `↑` / `↓` | Anywhere | Navigate |
| `Enter` | Anywhere | Select |
| `Space` | Now Playing | Pause / resume |
| `N` | Now Playing | Next track |
| `L` | Now Playing | Cycle repeat: off → one → all |
| `P` | Now Playing | Open Live Queue |
| `O` | Now Playing / Queue | Open Personal Playlist |
| `A` | Playlist views | Search and add a track |
| `D` | Playlist views | Remove the selected track |
| `C` | Personal Playlist | Clear the whole personal playlist |
| `T` | Anywhere | Theme selector |
| `V` | Anywhere | Cycle visualizer |
| `U` | Home | Install available update |
| `Q` / `Esc` | Anywhere | Back / quit |

### Personal Playlist vs Live Queue

Two separate lists, with a deliberate priority rule:

- **Personal Playlist** — tracks you explicitly added. Strictly FIFO, and it
  **always plays before** anything auto-generated. Lives in its own window
  (`[O]`). Adding to it never disturbs the currently playing track, and clearing
  it (`[C]`) never touches the Live Queue or playback.
- **Live Queue** — search results, autonomous mood-based play, and
  endless-listening refills. Only takes over once the Personal Playlist is empty.

When the Live Queue is exhausted and the Personal Playlist is empty, Moodwave
searches for similar music and refills automatically (endless listening).

---

## Environment Variables

Override any setting without editing the config file:

| Variable | Description |
|----------|-------------|
| `MOODWAVE_THEME` | Visual theme ID |
| `MOODWAVE_VISUAL` | Visual mode |
| `MOODWAVE_NO_COLOR=1` | Disable color |
| `MOODWAVE_NO_ANIMATION=1` | Disable animation |
| `MOODWAVE_NO_UNICODE=1` | Force ASCII output |
| `MOODWAVE_BACKEND` | Force audio backend |
| `MOODWAVE_VOLUME` | Volume (0–100) |
| `MOODWAVE_DEBUG=1` | Enable debug logging |
| `MOODWAVE_COLORS=0` | Disable color (alias) |
| `NO_COLOR` | Standard no-color flag |
| `NO_MOTION=1` | Prefer reduced motion |
| `JAMENDO_CLIENT_ID` | Enable Jamendo source |
| `LISTENBRAINZ_TOKEN` | Enable ListenBrainz source |
| `LISTENBRAINZ_USERNAME` | ListenBrainz username |
