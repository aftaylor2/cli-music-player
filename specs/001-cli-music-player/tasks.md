# Tasks: CLI Music Player

**Input**: Design documents from `/specs/001-cli-music-player/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Included per constitution principle III (Test-First Development). Test tasks target core logic packages (library, audio). TUI components are validated via end-to-end quickstart verification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project directory structure per plan: cmd/player/, internal/audio/, internal/library/, internal/tui/
- [x] T002 Initialize Go module (`go mod init`) and install dependencies: gopxl/beep v2, dhowden/tag, charmbracelet/bubbletea, charmbracelet/lipgloss, evertras/bubble-table, charmbracelet/bubbles
- [x] T003 [P] Add golangci-lint configuration file at .golangci.yml with go vet and staticcheck enabled

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data types, metadata parsing, directory scanning, and audio engine that ALL user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 [P] Define Track struct and AudioFormat enum (MP3, OGG) in internal/library/track.go per data-model.md
- [x] T005 [P] Define PlayStatus enum (Stopped, Playing, Paused) and PlayerState struct in internal/audio/player.go per data-model.md
- [x] T006 Write table-driven tests for metadata parsing (ID3 tags, Vorbis comments, missing tags fallback) in internal/library/track_test.go
- [x] T007 Implement metadata parsing using dhowden/tag in internal/library/track.go — NewTrackFromFile(path) returns Track with extracted title, artist, album, genre, track number; fallback to filename for missing title, "Unknown" for other missing fields
- [x] T008 Write table-driven tests for directory scanner (recursive scan, skip non-audio files, empty dir) in internal/library/scanner_test.go
- [x] T009 Implement recursive directory scanner in internal/library/scanner.go — ScanDirectory(root string) returns []Track; walks directory tree, filters by .mp3/.ogg extension (case-insensitive), calls NewTrackFromFile for each match
- [x] T010 [P] Implement audio decoder in internal/audio/decoder.go — OpenDecoder(path string, format AudioFormat) returns beep.StreamSeekCloser; uses beep/mp3 for MP3 and beep/vorbis for OGG
- [x] T011 Implement audio player engine in internal/audio/player.go — Player struct with Play(track), Pause(), Resume(), Stop(), IsPlaying(), Elapsed() methods; uses beep.Speaker for output, beep.Ctrl for pause/resume; initializes speaker with beep.Resample for sample rate normalization

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 - Browse and Play Music from a Directory (Priority: P1) MVP

**Goal**: User launches player with a directory, sees a table of tracks with metadata, selects a track to play, sees playback progress. Auto-advances to next track when current track ends.

**Independent Test**: Run `./player /path/to/music`, verify table displays, arrow keys navigate, Enter plays a track, Space pauses/resumes, playback position updates, auto-advances on track end.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T012 [P] [US1] Write tests for Library struct (NewLibrary, track count, empty library) in internal/library/library_test.go

### Implementation for User Story 1

- [x] T013 [US1] Implement Library struct in internal/library/library.go — NewLibrary(rootDir string) scans directory via ScanDirectory, stores []Track; Tracks() returns the full list; Len() returns count
- [x] T014 [P] [US1] Define key bindings in internal/tui/keys.go — keymap struct with bindings for: up/down/pgup/pgdown navigation, Enter (play), Space (pause/resume), s (stop), n (next), p (previous), q (quit) using bubbles/key
- [x] T015 [US1] Implement table view in internal/tui/table.go — configure evertras/bubble-table with columns: #, Title, Artist, Album, Genre, Duration; populate rows from []Track; handle row selection; truncate long text with ellipsis
- [x] T016 [P] [US1] Implement playback status bar in internal/tui/controls.go — render current track info (title - artist), play/pause/stop status icon, elapsed/total duration progress using lipgloss styling
- [x] T017 [US1] Implement main TUI model (Init/Update/View) in internal/tui/model.go — Model struct composing table, controls, audio player, library; handle key messages to dispatch play/pause/stop/next/prev; tick command to update elapsed time display every 500ms
- [x] T018 [US1] Add auto-advance logic in internal/tui/model.go — when audio player reports track ended (via tick check), advance to next track in visible list; stop after last track per FR-011
- [x] T019 [US1] Implement CLI entry point in cmd/player/main.go — parse os.Args for optional directory (default "."), validate directory exists (exit 1) and is directory (exit 2), create Library, create and run bubbletea.Program with TUI model

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Search Tracks (Priority: P2)

**Goal**: User presses `/` to enter search mode, types a query, and the table filters in real time to matching tracks across title, artist, and album fields.

**Independent Test**: With US1 complete, press `/`, type a query, verify table filters. Clear input or press Escape to restore full list.

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T020 [P] [US2] Write tests for Library.Search (case-insensitive match, no results, empty query returns all, matches across title/artist/album) in internal/library/library_test.go

### Implementation for User Story 2

- [x] T021 [US2] Add Search(query string) []Track method to Library in internal/library/library.go — case-insensitive substring match across Title, Artist, and Album fields; empty query returns all tracks
- [x] T022 [US2] Implement search input component in internal/tui/search.go — wraps bubbles/textinput; activated by `/` key; Escape dismisses and clears; renders search bar with lipgloss styling above or below table
- [x] T023 [US2] Integrate search mode into TUI model in internal/tui/model.go — when search active, filter table rows via Library.Search on each keystroke; current playing track continues but next/prev navigate within filtered list per FR-012

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Filter by Artist, Album, or Genre (Priority: P3)

**Goal**: User presses number keys (1-4) to switch between song/artist/album/genre views. Selecting a group drills down to show its tracks. Escape/Backspace returns to the previous level.

**Independent Test**: With US1 complete, press `2` for artist view, select an artist, see their tracks, play one. Press Escape to go back.

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T024 [P] [US3] Write tests for Library.GroupByArtist, GroupByAlbum, GroupByGenre, FilterByArtist, FilterByAlbum in internal/library/library_test.go

### Implementation for User Story 3

- [x] T025 [US3] Add GroupByArtist/GroupByAlbum/GroupByGenre and FilterByArtist/FilterByAlbum methods to Library in internal/library/library.go — GroupBy returns map[string][]Track; FilterBy returns []Track for a given key
- [x] T026 [US3] Implement view mode switching in internal/tui/views.go — ViewMode enum (Songs, Artists, Albums, Genres); render grouped list (show group names as rows) or track list depending on mode and drill-down state; manage navigation stack for back navigation
- [x] T027 [US3] Integrate view modes into TUI model in internal/tui/model.go — keybindings 1/2/3/4 switch ViewMode; Enter on a group row drills down; Escape/Backspace pops navigation stack; visible list updates for playback per FR-012

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Edge case handling, code quality, and end-to-end validation

- [x] T028 [P] Handle edge cases across packages: nonexistent directory error message (cmd/player/main.go), graceful skip on undecodable files (internal/library/scanner.go), handle deleted file during playback by advancing to next track (internal/audio/player.go)
- [x] T029 [P] Code quality pass: run gofmt on all files, verify go vet and golangci-lint produce zero warnings, fix any issues
- [ ] T030 Run quickstart.md end-to-end validation: build binary, launch with test directory containing MP3 and OGG files, verify all controls from quickstart.md Controls table work correctly

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in priority order (P1 → P2 → P3)
  - US2 and US3 build on US1's TUI model but are independently testable
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) — No dependencies on other stories
- **User Story 2 (P2)**: Depends on US1 TUI model (T017) being complete — adds search overlay
- **User Story 3 (P3)**: Depends on US1 TUI model (T017) being complete — adds view mode switching

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Data layer (Library methods) before UI components
- UI components before TUI model integration
- Story complete before moving to next priority

### Parallel Opportunities

- T004 and T005 can run in parallel (different packages)
- T006/T007 (metadata) and T008/T009 (scanner) can run in parallel after T004
- T010 (decoder) can run in parallel with library tasks
- T014 (keys) and T016 (controls) can run in parallel with T015 (table)
- T012 (library tests) can run in parallel with T014 and T016
- T020 (search tests) and T024 (grouping tests) are independent
- T028 and T029 can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch parallelizable tasks together:
Task: "T012 [P] [US1] Write tests for Library struct in internal/library/library_test.go"
Task: "T014 [P] [US1] Define key bindings in internal/tui/keys.go"
Task: "T016 [P] [US1] Implement playback status bar in internal/tui/controls.go"

# Then sequential (depends on above):
Task: "T013 [US1] Implement Library struct in internal/library/library.go"
Task: "T015 [US1] Implement table view in internal/tui/table.go"
Task: "T017 [US1] Implement main TUI model in internal/tui/model.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Build and run `./player /path/to/music`
5. Verify: table displays, playback works, auto-advance works

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → MVP!
3. Add User Story 2 → Test search independently
4. Add User Story 3 → Test view modes independently
5. Polish → Final validation with quickstart.md

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Audio player runs in a separate goroutine; communicate with TUI via Bubble Tea Cmd/Msg pattern
- beep.Speaker.Init must be called once; use beep.Resample for sample rate normalization
