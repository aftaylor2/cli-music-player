# Research: CLI Music Player

**Feature Branch**: `001-cli-music-player`
**Date**: 2026-03-17

## Audio Playback Library

**Decision**: gopxl/beep v2

**Rationale**: gopxl/beep is the actively maintained successor to
faiface/beep. It provides built-in MP3 and OGG decoding, a
composable Streamer interface for effects/mixing, and built-in
pause/resume/seek via Ctrl and Loop2 wrappers. It uses
ebitengine/oto v3 under the hood for cross-platform audio output.

**Alternatives considered**:
- **ebitengine/oto directly**: Lower-level, would require manually
  integrating separate MP3/OGG decoders. More work for no benefit
  in this use case.
- **hajimehoshi/oto**: Deprecated in favor of ebitengine/oto.

**Caveats**:
- Requires CGo (ALSA headers on Linux: `alsa-lib` on Arch,
  `libasound2-dev` on Debian).
- The underlying go-mp3 decoder is unmaintained but functional.
  Monitor for edge-case MP3 issues.

## Metadata / Tag Reading

**Decision**: dhowden/tag

**Rationale**: Provides a single unified interface for reading ID3v1,
ID3v2 (2.2/2.3/2.4), OGG Vorbis comments, and FLAC metadata.
Read-only API is simple: `tag.ReadFrom(reader)` returns a Tag
interface with Title(), Artist(), Album(), Genre(), etc. Actively
maintained (last activity Jan 2025).

**Alternatives considered**:
- **bogem/id3v2**: MP3-only, geared toward tag writing/editing.
  Overkill for read-only use and doesn't cover OGG.

## TUI Table Component

**Decision**: evertras/bubble-table

**Rationale**: Purpose-built table component for Bubble Tea with
flexible column sizing, row selection, filtering, sorting, and
pagination. Supports per-column and per-row styling via Lip Gloss.
The charmbracelet/bubbles table is simpler but lacks the filtering
and sorting features needed for a music player.

**Alternatives considered**:
- **charmbracelet/bubbles table**: Too basic — limited column
  configuration, no built-in filtering or sorting.
- **Custom implementation**: Unnecessary given bubble-table covers
  the requirements.

## Platform & Build Considerations

- **Linux audio**: ALSA is required. PulseAudio/PipeWire work via
  the ALSA compatibility layer.
- **Cross-compilation**: CGo complicates cross-compilation. Build
  on target platform or use Docker-based xgo for cross-builds.
- **Threading**: Beep/Oto has thread-safety requirements. Audio
  control must use Bubble Tea's message-passing (Cmd/Msg pattern)
  to communicate between the TUI goroutine and the audio goroutine.
- **Sample rates**: Test with both 44.1kHz and 48kHz content. Beep
  provides a Resample streamer for rate conversion.

## Duration Extraction

**Decision**: Use beep's MP3/Vorbis decoders for duration (they
expose stream length). For faster scanning, dhowden/tag does not
provide duration — so we decode each file briefly to get the sample
count and compute duration from sample rate. Cache results in memory
after initial scan.

## Final Dependency Stack

| Purpose          | Package                      |
|------------------|------------------------------|
| Audio playback   | gopxl/beep v2                |
| Audio output     | ebitengine/oto v3 (via beep) |
| MP3 decoding     | gopxl/beep/mp3               |
| OGG decoding     | gopxl/beep/vorbis            |
| Metadata tags    | dhowden/tag                  |
| TUI framework    | charmbracelet/bubbletea      |
| TUI styling      | charmbracelet/lipgloss       |
| Table component  | evertras/bubble-table        |
| Input/viewport   | charmbracelet/bubbles        |
