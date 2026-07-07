# Moodwave CLI — Status

**Current phase: MVP Foundation**
**Build status: ✓ Compiles and runs**
**Last updated: 2026-07-06**

---

## What works

| Feature | Status |
|---------|--------|
| `moodwave init` | ✅ Initializes config and cache directories |
| `moodwave scan` | ✅ Full repository analysis, mood inference |
| `moodwave mood` | ✅ Displays mood profile, scores, signals |
| `moodwave play` | ✅ Recommends tracks, starts audio backend |
| `moodwave status` | ✅ Shows system and session state |
| `moodwave config` | ✅ Shows config, points to config file |
| `moodwave theme` | ✅ Lists themes, switches runtime theme |
| `moodwave visual` | ✅ Lists modes, switches visual mode |
| `moodwave source` | ✅ Lists sources, runs health checks |
| `moodwave doctor` | ✅ Full system diagnostics |
| `moodwave queue` | ✅ Shows recommendation queue |
| `moodwave next` | ✅ Advances queue index |
| `moodwave stop` | ✅ Updates session playback state |
| `moodwave --help` | ✅ Full usage output |
| `moodwave --version` | ✅ Version and build time |

## Music sources

| Source | Status |
|--------|--------|
| Radio Browser | ✅ Implemented, queried first |
| LRCLIB lyrics | ✅ Implemented |
| MusicBrainz | ✅ Implemented (rate-limited, metadata only) |
| Jamendo | ✅ Implemented (optional, needs client_id) |
| ListenBrainz | 🔲 Stub — planned for Phase 2 |

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

## Mood engine — scoring rules

| Rule | Status |
|------|--------|
| Test coverage signal | ✅ |
| TODO/FIXME density | ✅ |
| Comment density | ✅ |
| Structure entropy | ✅ |
| Naming consistency | ✅ |
| Documentation | ✅ |
| Build system | ✅ |
| Dependency weight | ✅ |
| Primary language hints | ✅ |
| Language diversity | ✅ |
| Git churn | ✅ |
| Project size | ✅ |
| CI presence | ✅ |

## Visual engine

| Feature | Status |
|---------|--------|
| Dot-matrix banner | ✅ |
| Wave animation | ✅ |
| Spectrum bars | ✅ |
| Pulse animation | ✅ |
| Minimal mode | ✅ |
| Quiet (no animation) | ✅ |
| ANSI color + monochrome fallback | ✅ |
| Unicode + ASCII fallback | ✅ |
| Narrow terminal adaptation | ✅ |
| Reduced-motion support | ✅ |

## Platform support

| Platform | Status |
|----------|--------|
| Windows (x86_64) | ✅ Binary compiles and runs |
| Windows (ARM64) | ✅ Cross-compile target |
| macOS (Apple Silicon) | ✅ Cross-compile target |
| macOS (Intel) | ✅ Cross-compile target |
| Linux (x86_64) | ✅ Cross-compile target |
| Linux (ARM64) | ✅ Cross-compile target |

## Audio backends

| Backend | Platform | Status |
|---------|----------|--------|
| mpv | All | ✅ Detected and used if available |
| ffplay (ffmpeg) | All | ✅ Fallback |
| afplay | macOS | ✅ Native fallback |
| PowerShell WMP | Windows | ✅ Last resort |
| None | Any | ✅ Info-only mode, shows stream URL |

## Known gaps / planned

- [ ] Real-time audio visualization (sync to audio data, not time-based)
- [ ] ListenBrainz scrobbling
- [ ] Per-project config override (`.moodwave.toml`)
- [ ] `moodwave watch` — background mood auto-update daemon
- [ ] Terminal keypress handling during playback (live keyboard controls)
- [ ] Lyric display during playback
- [ ] Website (Phase 3)

---

*The CLI is the product. Everything else is secondary.*
