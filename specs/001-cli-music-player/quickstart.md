# Quickstart: CLI Music Player

## Prerequisites

- Go 1.22+ installed
- ALSA development headers (for audio playback):
  - Arch: `sudo pacman -S alsa-lib`
  - Debian/Ubuntu: `sudo apt install libasound2-dev`
  - Fedora: `sudo dnf install alsa-lib-devel`

## Build

```bash
cd /home/ataylor/Projects/cli-music-player
go build -o player ./cmd/player
```

## Run

```bash
# Play music from a specific directory
./player /path/to/music

# Play music from the current directory
./player
```

## Controls

| Key          | Action                    |
|--------------|---------------------------|
| ↑/↓          | Navigate track list       |
| Enter        | Play selected track       |
| Space        | Pause / Resume            |
| s             | Stop playback             |
| n             | Next track                |
| p             | Previous track            |
| /             | Enter search mode         |
| Escape        | Exit search / go back     |
| 1             | Song view (default)       |
| 2             | Artist view               |
| 3             | Album view                |
| 4             | Genre view                |
| q             | Quit                      |

## Verify

1. Prepare a directory with at least one `.mp3` and one `.ogg` file.
2. Run `./player /path/to/that/directory`.
3. Confirm the table displays track metadata (title, artist, album,
   genre, duration).
4. Select a track and press Enter — audio should play.
5. Press Space to pause, Space again to resume.
6. Press `/`, type a search query, confirm filtering works.
7. Press `q` to quit.

## Development

```bash
# Run tests
go test ./...

# Run with verbose test output
go test -v ./...

# Lint
golangci-lint run

# Format
gofmt -w .
```
