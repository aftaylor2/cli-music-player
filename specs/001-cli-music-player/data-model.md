# Data Model: CLI Music Player

**Feature Branch**: `001-cli-music-player`
**Date**: 2026-03-17

## Entities

### Track

Represents a single audio file with its metadata.

| Field       | Type            | Description                              |
|-------------|-----------------|------------------------------------------|
| FilePath    | string          | Absolute path to the audio file          |
| Title       | string          | Track title (from tags, fallback: filename) |
| Artist      | string          | Artist name (fallback: "Unknown")        |
| Album       | string          | Album name (fallback: "Unknown")         |
| Genre       | string          | Genre (fallback: "Unknown")              |
| TrackNumber | int             | Track number within album (0 if missing) |
| Duration    | time.Duration   | Total track duration                     |
| Format      | AudioFormat     | Enum: MP3, OGG                           |

**Validation rules**:
- FilePath MUST be non-empty and point to a readable file.
- Duration MUST be >= 0.
- Title falls back to the filename (without extension) when tag is
  missing or empty.

### AudioFormat

Enum type representing supported audio formats.

| Value | Description       |
|-------|-------------------|
| MP3   | MPEG Audio Layer 3 |
| OGG   | Ogg Vorbis        |

### Library

The in-memory collection of all scanned tracks. Not persisted.

| Field    | Type      | Description                              |
|----------|-----------|------------------------------------------|
| Tracks   | []Track   | All tracks found during directory scan   |
| RootDir  | string    | The root directory that was scanned      |

**Operations**:
- `Search(query string) []Track` — case-insensitive substring match
  across title, artist, and album fields.
- `GroupByArtist() map[string][]Track`
- `GroupByAlbum() map[string][]Track`
- `GroupByGenre() map[string][]Track`
- `FilterByArtist(artist string) []Track`
- `FilterByAlbum(album string) []Track`

### PlayerState

Represents the current playback state. Driven by the audio engine,
reflected in the TUI.

| Field       | Type          | Description                          |
|-------------|---------------|--------------------------------------|
| Status      | PlayStatus    | Playing, Paused, or Stopped          |
| CurrentTrack| *Track        | Pointer to currently loaded track    |
| Elapsed     | time.Duration | Current playback position            |
| Volume      | float64       | Volume level 0.0–1.0 (future use)   |

### PlayStatus

Enum for playback state.

| Value   | Description                       |
|---------|-----------------------------------|
| Stopped | No track loaded or playback ended |
| Playing | Audio is actively playing         |
| Paused  | Playback is paused                |

## Relationships

```text
Library 1──* Track
PlayerState *──1 Track (current)
```

- A Library contains zero or more Tracks.
- PlayerState references at most one Track at a time.

## State Transitions (PlayerState)

```text
Stopped ──[select track]──► Playing
Playing ──[pause]──────────► Paused
Playing ──[stop]───────────► Stopped
Playing ──[track ends]─────► Stopped
Paused  ──[resume]─────────► Playing
Paused  ──[stop]───────────► Stopped
Any     ──[select track]───► Playing (restarts with new track)
```
