<!--
  Sync Impact Report
  ───────────────────
  Version change: 0.0.0 → 1.0.0 (initial ratification)
  Modified principles: N/A (initial version)
  Added sections:
    - Core Principles (5 principles)
    - Technology Stack
    - Development Workflow
    - Governance
  Removed sections: N/A
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ no changes needed (generic)
    - .specify/templates/spec-template.md ✅ no changes needed (generic)
    - .specify/templates/tasks-template.md ✅ no changes needed (generic)
    - .specify/templates/commands/*.md ✅ no command files present
  Follow-up TODOs: none
-->

# CLI Music Player Constitution

## Core Principles

### I. TUI-First Design

All user interaction MUST occur through a terminal user interface
built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).
Visual styling MUST use
[Lip Gloss](https://github.com/charmbracelet/lipgloss) where needed.

- The TUI MUST present a table-like display showing: song title,
  artist, album, genre, track duration, and playback position.
- The interface MUST support keyboard-driven navigation, search,
  and filtering by artist, album, or song.
- The player MUST accept a directory path as a CLI argument or
  default to the current working directory when none is provided.

### II. Clean Package Architecture

The codebase MUST be organized into well-defined Go packages with
clear responsibilities:

- **Separation of concerns**: audio decoding, metadata parsing,
  TUI rendering, and playback control MUST reside in distinct
  packages.
- No package may import a sibling package that creates a circular
  dependency.
- Public APIs between packages MUST be minimal and intentional.

### III. Test-First Development

Tests MUST be written before implementation code (red-green-refactor).

- Unit tests MUST cover metadata parsing, playlist logic, and
  filtering/search operations.
- Integration tests MUST verify end-to-end playback initialization
  and TUI model updates.
- Table-driven tests are preferred for Go test files.

### IV. Audio Format Support

The player MUST support MP3 and OGG (Vorbis) file formats.

- Unsupported files in a scanned directory MUST be silently skipped
  without errors shown to the user.
- Metadata extraction (ID3 for MP3, Vorbis comments for OGG) MUST
  be handled gracefully — missing tags MUST fall back to filename
  or "Unknown" values.
- Adding future format support MUST NOT require changes to the TUI
  or playlist packages (open/closed principle at the decoder boundary).

### V. Simplicity

Start with the minimum viable feature set; avoid premature abstraction.

- YAGNI: do not build plugin systems, configuration files, or
  network features unless explicitly requested.
- Prefer standard library and Charm ecosystem packages over
  third-party dependencies.
- A single binary with zero required configuration MUST be the
  deployment target.

## Technology Stack

- **Language**: Go (latest stable, currently 1.22+)
- **TUI framework**: Bubble Tea (charmbracelet/bubbletea)
- **Styling**: Lip Gloss (charmbracelet/lipgloss)
- **Audio playback**: To be determined during planning (e.g.,
  faiface/beep, hajimehoshi/oto, or equivalent)
- **Metadata**: Tag-reading library supporting ID3v2 and Vorbis
  comments (e.g., dhowden/tag)
- **Build**: `go build` / `go install`; no external build tools
  required
- **Testing**: `go test` with standard testing package

## Development Workflow

- All changes MUST compile with `go vet` and pass `staticcheck`
  (or golangci-lint) with zero warnings before merge.
- Each user story MUST be deliverable and testable independently.
- Commits MUST be atomic — one logical change per commit.
- Code formatting MUST use `gofmt` (enforced, not optional).

## Governance

This constitution is the authoritative source of project principles.
All design decisions, code reviews, and implementation plans MUST
be consistent with these principles.

- **Amendments**: Any change to this constitution MUST be documented
  with a version bump, a rationale, and a review of dependent
  templates for consistency.
- **Versioning**: MAJOR.MINOR.PATCH semantic versioning. MAJOR for
  principle removals or incompatible redefinitions, MINOR for new
  principles or material expansions, PATCH for clarifications.
- **Compliance**: Plan and spec documents MUST reference the
  constitution version they were authored against.

**Version**: 1.0.0 | **Ratified**: 2026-03-17 | **Last Amended**: 2026-03-17
