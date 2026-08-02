# Moodwave CLI — Status

**Current phase: v2.0.0 — Terminal Rebuild**
**Build status: ✓ Compiles, vets, and tests clean on Windows / macOS / Linux**
**Last updated: 2026-08-03**

---

## Commands

| Feature | Status |
|---------|--------|
| `moodwave` (interactive TUI) | ✅ Full Bubble Tea TUI — Home, Search, Now Playing, Live Queue, Personal Playlist, Theme Select, Doctor |
| `moodwave init` | ✅ Initializes config and cache directories |
| `moodwave scan` | ✅ Full repository analysis, mood inference |
| `moodwave mood` | ✅ Displays mood profile, scores, signals |
| `moodwave play` | ✅ Recommends tracks, starts audio backend |
| `moodwave search` | ✅ Live YouTube search and selection |
| `moodwave status` | ✅ Shows system and session state |
| `moodwave config` | ✅ Shows config, points to config file |
| `moodwave theme` | ✅ Lists themes, switches runtime theme |
| `moodwave visual` | ✅ Lists modes, switches visual mode |
| `moodwave source` | ✅ Lists sources, runs health checks |
| `moodwave doctor` | ✅ Interactive diagnostics TUI with one-key auto-fix |
| `moodwave queue` | ✅ Shows recommendation queue |
| `moodwave next` | ✅ Advances queue index |
| `moodwave stop` | ✅ Updates session playback state |
| `moodwave update` | ✅ In-place self-update from GitHub releases |
| `moodwave --help` | ✅ Full usage output |
| `moodwave --version` | ✅ Version and build time |

## Playback engine

| Feature | Status |
|---------|--------|
| Single-owner `PlaybackQueue` (no dual-slice desync) | ✅ |
| Personal Playlist — strict FIFO, always plays first | ✅ |
| One consolidated `resolveOutcome()` next-track decision | ✅ |
| Repeat modes: off / one / all | ✅ |
| Manual skip bypasses repeat-one | ✅ |
| Consecutive-failure cap (force-advance past broken tracks) | ✅ |
| Controller-tagged messages (discards stale/double-audio races) | ✅ |
| Endless listening — similar-music refill on queue exhaustion | ✅ |
| Resolve retry-once on transient network failure | ✅ |

## Music sources

| Source | Status |
|--------|--------|
| YouTube (yt-dlp + Invidious fallback) | ✅ Implemented, primary source |
| Radio Browser | ✅ Implemented |
| LRCLIB lyrics | ✅ Implemented |
| MusicBrainz | ✅ Implemented (rate-limited, metadata only) |
| Jamendo | ✅ Implemented (optional, needs client_id) |
| ListenBrainz | 🔲 Stub — planned |

## Scanner signals extracted

| Signal | Status |
|--------|--------|
| Language detection (40+ languages) | ✅ |
| Build system detection | ✅ |
| Dependency manifest detection | ✅ |
| Test file ratio | ✅ |
| TODO/FIXME density | ✅ |
| Comment density | ✅ |
| Directory structure entropy | ✅ |
| Naming consistency | ✅ |
| Documentation presence | ✅ |
| CI configuration | ✅ |
| Git branch / churn / age | ✅ |
| Semantic vocabulary counts | ✅ |

## Mood engine — scoring rules

All 13 heuristic rules implemented: test coverage, TODO/FIXME density, comment
density, structure entropy, naming consistency, documentation, build system,
dependency weight, primary language hints, language diversity, git churn,
project size, and CI presence. ✅

## Visual engine

| Feature | Status |
|---------|--------|
| Bubble Tea TUI (async, non-blocking) | ✅ |
| 24 terminal visualizers | ✅ |
| Doom-style cellular-automaton fire | ✅ |
| Dynamic sizing to terminal dimensions | ✅ |
| 7 TUI themes + live preview selector | ✅ |
| Hand-drawn animated pixel pet (Home / Search) | ✅ |
| ASCII banner renderer | ✅ |
| ANSI color + monochrome fallback | ✅ |
| Unicode + ASCII fallback | ✅ |
| Narrow terminal adaptation | ✅ |
| Reduced-motion support | ✅ |

## Self-update

| Feature | Status |
|---------|--------|
| GitHub release check (`internal/updater`) | ✅ |
| Automatic check on every TUI launch | ✅ |
| Home-screen banner + `[U]` install | ✅ |
| `moodwave update` command (shared implementation) | ✅ |
| Atomic binary swap with rollback on failure | ✅ |
| Silent failure on background check (never interrupts) | ✅ |

## Platform support

| Platform | Status |
|----------|--------|
| Windows (x86_64) | ✅ Builds and runs |
| Windows (ARM64) | ✅ Cross-compile verified |
| macOS (Apple Silicon) | ✅ Cross-compile verified |
| macOS (Intel) | ✅ Cross-compile verified |
| Linux (x86_64) | ✅ Cross-compile verified |
| Linux (ARM64) | ✅ Cross-compile verified |
| Linux (ARMv7) | ✅ Cross-compile verified |

## Audio backends

| Backend | Platform | Status |
|---------|----------|--------|
| mpv | All | ✅ Detected and used if available |
| ffplay (ffmpeg) | All | ✅ Fallback + auto-download if missing |
| afplay | macOS | ✅ Native CoreAudio fallback |
| Media Foundation | Windows | ✅ Native backend |
| None | Any | ✅ Info-only mode, shows stream URL |

## CI / release

| Job | Status |
|-----|--------|
| CLI build + vet + test (Windows / macOS / Linux) | ✅ Passing |
| Cross-compile all 7 targets | ✅ Passing |
| Website type-check + build | ✅ Passing |
| Tagged release → 7 platform binaries + checksums | ✅ Configured |

> Tests that hit live third-party APIs (YouTube, Jamendo) are gated behind
> `testing.Short()` so CI runs deterministically with `go test ./... -short`.
> Run the full suite locally without `-short` to exercise them.

## Known gaps / planned

- [ ] Real-time audio visualization (sync to audio data, not time-based)
- [ ] ListenBrainz scrobbling
- [ ] `moodwave watch` — background mood auto-update daemon
- [ ] Lyric display during playback
- [ ] Persisting the Personal Playlist across sessions
- [ ] Windows Terminal redraw artifact on certain rapid view transitions

---

*The CLI is the product. Everything else is secondary.*
