# Feature Specification: CLI Music Player

**Feature Branch**: `001-cli-music-player`
**Created**: 2026-03-17
**Status**: Draft
**Input**: User description: "Creating a CLI / TUI music player (mp3 / ogg) using go. For the TUI use bubbletea and if needed lipgloss for styling. we want to be able to open in a dir or specify the dir. be able to display the songs, artist, album, genre, song time, duration, etc in a table like display. allow for searching. select by artist/album/song."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse and Play Music from a Directory (Priority: P1)

As a user, I want to launch the music player pointing at a directory
of MP3/OGG files and see all tracks listed in a table with metadata
(title, artist, album, genre, duration). I can navigate the list with
keyboard controls and press Enter to play a selected track. Playback
status (current position, total duration) is shown in the UI.

**Why this priority**: This is the core value proposition — without
browsing and playing music, there is no product.

**Independent Test**: Launch the binary with a directory containing
MP3 and OGG files. Verify the table displays metadata, arrow keys
navigate, Enter starts playback, and the progress/duration updates.

**Acceptance Scenarios**:

1. **Given** a directory with MP3 and OGG files, **When** I run
   `player /path/to/music`, **Then** I see a table listing all
   supported audio files with title, artist, album, genre, and
   duration columns.
2. **Given** the track list is displayed, **When** I select a track
   and press Enter, **Then** playback begins and the UI shows the
   current playback position and total duration.
3. **Given** no directory argument is provided, **When** I run
   `player`, **Then** the current working directory is scanned for
   audio files.
4. **Given** a directory with no supported audio files, **When** I
   launch the player, **Then** the UI displays an informative
   "No audio files found" message.

---

### User Story 2 - Search Tracks (Priority: P2)

As a user, I want to type a search query to filter the track list
in real time so I can quickly find a specific song, artist, or album
without scrolling through hundreds of entries.

**Why this priority**: Search makes the player usable for large
music collections. It builds on the table display from US1.

**Independent Test**: With US1 complete, press `/` to enter search
mode, type a query, and verify the table filters to matching rows.

**Acceptance Scenarios**:

1. **Given** the track list is displayed, **When** I press `/` and
   type "Beatles", **Then** only tracks whose title, artist, or
   album contain "Beatles" (case-insensitive) are shown.
2. **Given** a search filter is active, **When** I clear the search
   input, **Then** all tracks are shown again.
3. **Given** a search filter is active, **When** I press Escape,
   **Then** the search is dismissed and the full list is restored.

---

### User Story 3 - Filter by Artist, Album, or Genre (Priority: P3)

As a user, I want to switch between views that group tracks by
artist, album, or genre so I can browse my music collection
hierarchically — select an artist, then see their albums, then
tracks.

**Why this priority**: Hierarchical browsing adds discoverability
but is not required for basic playback or search.

**Independent Test**: With US1 complete, use a keybinding to switch
to "artist view", select an artist, see their tracks, and play one.

**Acceptance Scenarios**:

1. **Given** the track list is displayed, **When** I press a
   keybinding (e.g., `1` for songs, `2` for artists, `3` for
   albums, `4` for genres), **Then** the view changes to group
   tracks accordingly.
2. **Given** I am in artist view, **When** I select an artist and
   press Enter, **Then** I see only that artist's tracks.
3. **Given** I am in a filtered view, **When** I press Backspace or
   Escape, **Then** I return to the previous grouping level.

---

### Edge Cases

- What happens when an audio file has corrupt or missing metadata
  tags? Display "Unknown" for missing fields; skip files that
  cannot be decoded at all.
- What happens when the specified directory does not exist? Exit
  with a clear error message and non-zero exit code.
- What happens when a file is deleted while the player is running?
  Skip to the next track gracefully if the currently playing file
  is removed.
- What happens with very long track titles or artist names? Truncate
  with ellipsis to fit the column width.
- What happens with nested subdirectories? Recursively scan for
  audio files.

## Clarifications

### Session 2026-03-17

- Q: When a track finishes playing, what should happen? → A: Auto-advance to the next track in the current list. Playback stops after the last track.
- Q: Is there a separate play queue, or does "play" always mean playing from the currently visible list? → A: Play from visible list only (no separate queue).
- Q: Should volume control be included in this version? → A: No, defer to system volume controls.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST scan a given directory (recursively) for
  MP3 and OGG files and ignore all other file types.
- **FR-002**: System MUST extract metadata (title, artist, album,
  genre, track number, duration) from each audio file.
- **FR-003**: System MUST display tracks in a table with columns:
  title, artist, album, genre, and duration.
- **FR-004**: System MUST allow keyboard navigation (up/down/
  page-up/page-down) through the track list.
- **FR-005**: System MUST play a selected track and display current
  playback position alongside total duration.
- **FR-006**: System MUST support play, pause, stop, next track,
  and previous track controls.
- **FR-011**: System MUST auto-advance to the next track in the
  current list when a track finishes. Playback stops after the
  last track in the list.
- **FR-012**: Playback MUST operate on the currently visible list.
  There is no separate play queue. When a search or filter changes
  the visible list, the current track continues playing but
  next/previous navigate within the new visible list.
- **FR-007**: System MUST accept a directory path as a CLI argument
  or default to the current working directory.
- **FR-008**: System MUST provide a search mode that filters the
  track list by matching against title, artist, and album fields
  (case-insensitive).
- **FR-009**: System MUST support grouping/filtering tracks by
  artist, album, or genre via keybindings.
- **FR-010**: System MUST handle missing metadata gracefully by
  displaying "Unknown" or the filename as a fallback.

### Key Entities

- **Track**: Represents a single audio file. Attributes: file path,
  title, artist, album, genre, track number, duration, format
  (mp3/ogg).
- **Library**: The collection of all tracks found in the scanned
  directory. Supports filtering, searching, and grouping operations.
- **Player State**: Current playback status including: playing/
  paused/stopped, current track, elapsed time. Volume control is
  deferred to the system mixer (out of scope for v1).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: User can launch the player, see a populated track list,
  and begin playback within 3 seconds for a library of up to 1000
  tracks.
- **SC-002**: Search results update within 100ms of keystroke for
  libraries up to 10,000 tracks.
- **SC-003**: The player binary is a single statically-compiled
  executable with zero runtime dependencies.
- **SC-004**: All supported audio formats (MP3, OGG) play without
  audible glitches or gaps on Linux.
