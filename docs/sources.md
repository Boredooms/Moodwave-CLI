# Moodwave CLI — API Sources and Integration Map

## 1. Goal

This document defines the external sources, APIs, and reference repositories that Moodwave CLI can use for music discovery, metadata, recommendations, lyrics, station fallback, and terminal playback inspiration.

The core rule is simple: **use documented, legal, stable APIs first**. Do not depend on scraping or unstable unofficial endpoints for the main product.

## 2. Source Strategy

Moodwave CLI should use a multi-source adapter design:

- one metadata source
- one user-history / recommendation source
- one lyrics source
- one radio/station fallback source
- one local playback backend
- optional model-based mood inference
- optional open-source visual renderer support

This keeps the CLI resilient when one provider fails.

## 3. Recommended API Sources

### 3.1 MusicBrainz
**Use for:** core metadata lookup, artist/release/recording identification, canonical music database references.

**Why it matters:** MusicBrainz is an open music encyclopedia with a REST API for music metadata. The API returns XML or JSON and is aimed at media players, CD rippers, taggers, and similar apps. It also requires a meaningful User-Agent and limits clients to one call per second. 

**Best use in Moodwave CLI:**
- resolve track and artist metadata
- canonicalize song names
- attach stable IDs to tracks
- build a metadata cache
- power downstream matching logic

**Usage pattern:**
1. Search artist/recording/release.
2. Normalize results into internal track objects.
3. Store only IDs and small metadata.
4. Rate-limit requests carefully.

### 3.2 ListenBrainz
**Use for:** listening history, recommendations, statistics, and user-profile-driven suggestions.

**Why it matters:** ListenBrainz publicly stores listens and exposes API endpoints for reading and submitting listens. The API uses HTTPS and user tokens for authentication. The docs also expose metadata endpoints that fetch MusicBrainz-linked recording details.

**Best use in Moodwave CLI:**
- learn user taste over time
- support personal recommendation history
- optionally submit listening sessions
- retrieve recording metadata tied to listens

**Usage pattern:**
1. Authenticate with a user token.
2. Fetch recent listens.
3. Convert them into taste vectors.
4. Blend history with current mood inference.

### 3.3 Jamendo
**Use for:** discoverable legal music catalogs, radio-style streams, music discovery, and track browsing.

**Why it matters:** Jamendo’s API advertises more than 20 read methods for a catalog of roughly half a million tracks, plus search, radios, OAuth2 authentication, and user library write methods. The radios endpoint can produce radio lists and stream information, and Jamendo notes that commercial use requires a commercial license.

**Best use in Moodwave CLI:**
- discover streamable tracks
- build genre/mood collections
- support radio-like playback fallback
- provide an easy legal source for MVP playback

**Usage pattern:**
1. Query catalog/search endpoints.
2. Rank tracks by mood fit.
3. Resolve a radio or stream option.
4. Keep compliance rules visible.

### 3.4 Radio Browser
**Use for:** open radio station fallback when track streaming is unavailable or when the user wants ambient continuous audio.

**Why it matters:** Radio Browser is a completely free and open-source API. It can be used directly or via third-party libraries, and the project explicitly allows self-hosting, forking, and integration in free or non-free software.

**Best use in Moodwave CLI:**
- fallback ambient playback
- genre-based radio modes
- coding-mood station presets
- lightweight no-library streaming path

**Usage pattern:**
1. Search stations by tags, country, or name.
2. Normalize station metadata.
3. Rank stations by mood match.
4. Stream one selected station.

### 3.5 LRCLIB
**Use for:** synchronized lyrics.

**Why it matters:** LRCLIB is an open-source collection of synchronized song lyrics, and the server is implemented in Rust with Axum and SQLite. It exposes machine-friendly APIs and is a good fit for terminal lyric overlays.

**Best use in Moodwave CLI:**
- optional lyrics display
- timed lyric syncing
- minimal UI overlays during playback

**Usage pattern:**
1. Fetch track match by title/artist.
2. Retrieve synchronized lyrics.
3. Render them in a low-distraction mode.

## 4. Reference Open-Source Repositories

These are not primary APIs, but they are important design references.

### 4.1 kew
A terminal music player that supports auto-generated playlists, gapless playback, private/offline usage, spectrum visualization, sixel cover art, and quick search. It is a strong reference for terminal-native music UX.

### 4.2 musikcube
A cross-platform terminal-based audio engine, library, player, and server written in C++. It runs on Windows, macOS, Linux, and Raspberry Pi, and it is a strong reference for a native, small-footprint playback architecture.

### 4.3 Spotube
A cross-platform open-source music streaming platform with a plugin-friendly design and a focus on being lightweight. It is useful as a reference for source separation, extensibility, and resource-conscious playback.

### 4.4 cava
A cross-platform audio visualizer for terminal or desktop. It is useful as a reference for spectrum-style rendering, audio-reactive motion, and portable visualization design.

## 5. Recommended Source Stack for Moodwave CLI

### Tier 1: Core metadata
- MusicBrainz

### Tier 2: Listening history and recommendation memory
- ListenBrainz

### Tier 3: Legal discoverable music
- Jamendo

### Tier 4: Ambient fallback
- Radio Browser

### Tier 5: Lyrics
- LRCLIB

### Tier 6: Playback references
- kew
- musikcube
- Spotube

### Tier 7: Visual references
- cava

## 6. Adapter Contract

Each source module should expose the same internal interface.

```text
search(query, filters)
resolve(id)
normalize(raw)
rank(candidate, context)
health_check()
stream(candidate)
fallback(candidate)
```

This makes it easy to swap providers without changing the CLI core.

## 7. Internal Data Model

### Track object
- id
- source
- title
- artist
- album
- duration
- bpm
- energy
- mood_tags
- stream_url
- preview_url
- lyrics_available
- cache_ttl

### Mood profile
- mood label
- confidence
- contributing signals
- preferred track traits
- transition policy

### Session profile
- repo path
- language mix
- git activity
- user overrides
- current playback source
- recent tracks

## 8. Mood-to-Music Matching Flow

1. Scan the repo.
2. Build a mood profile.
3. Expand that mood into track traits.
4. Query metadata sources.
5. Rank candidates.
6. Select a stream or station.
7. Start playback.
8. Render terminal visuals.
9. Re-evaluate on meaningful code changes.

## 9. Lightweight Model Strategy

The CLI should not require a heavy model to work.

### Recommended order
1. rules and heuristics
2. small embedding model
3. compact classifier
4. hybrid confidence scoring

### Use cases for models
- codebase tone detection
- session-style inference
- music trait matching
- transition suggestions

## 10. Streaming and Storage Rules

- do not permanently download media by default
- do not keep full tracks in memory
- cache only metadata and small previews where legal
- prefer stream URLs or station streams
- free space after playback stops
- keep local state compact

## 11. Terminal Rendering Sources

The visual system can be built from scratch, but it may also borrow ideas from open-source visualizer projects like cava.

### Visual primitives
- dot matrix text
- wave lines
- spectrum bars
- pulse rings
- ambient background noise
- ASCII art scenes
- progress meters
- mood badges

### Fallback rules
- no color -> monochrome
- no animation -> static frames
- no unicode -> ASCII only
- no TTY -> simple log output

## 12. Platform Support Notes

### Windows
- PowerShell bootstrap
- batch wrapper support
- native console-friendly rendering
- minimal dependency chain

### macOS
- shell install script
- native audio backend
- ANSI and Unicode visuals

### Linux
- curl install path
- distro-friendly packaging later
- terminal-first playback and visuals

## 13. What to Avoid

- scraping music sites as the main source
- using undocumented endpoints as a core dependency
- depending on one provider only
- heavy media caches
- large model dependencies in the default path
- browser-only implementations for the CLI core

## 14. Suggested MVP Source Set

If the first release must stay small, use:
- MusicBrainz
- ListenBrainz
- Jamendo
- Radio Browser
- LRCLIB

Then study:
- kew
- musikcube
- Spotube
- cava

## 15. Later Expansion Ideas

After the MVP, the source layer can expand into:
- additional station catalogs
- more metadata mirrors
- local library scanning
- personalized recommendation layers
- theme packs
- alternate terminal visual engines

## 16. Definition of a Good Source Layer

A good source layer is one that:
- keeps the CLI legal and stable
- gives the recommender enough metadata
- supports fallback modes
- works on multiple platforms
- avoids unnecessary storage
- can be extended later without rewriting the core
