# CLI Music Player

A terminal-based music player for MP3, OGG, and FLAC files, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Browse your music library in an interactive table, search and filter by artist, album, or genre, and control playback — all from the terminal.

![Player](https://storage.googleapis.com/ataylor-public/player/player.png)

![Album Art & Metadata](https://storage.googleapis.com/ataylor-public/player/player-album-artwork.png)

## Prerequisites

- Go 1.22+
- **macOS**: Works out of the box (uses CoreAudio)
- **Linux**: ALSA development headers:
  - **Arch**: `sudo pacman -S alsa-lib`
  - **Debian/Ubuntu**: `sudo apt install libasound2-dev`
  - **Fedora**: `sudo dnf install alsa-lib-devel`

## Install

```bash
git clone https://github.com/aftaylor2/cli-music-player.git
cd cli-music-player
make build
make install
```

## Usage

```bash
# Play music from a directory
./player /path/to/music

# Play music from the current directory
./player
```

## Controls

| Key       | Action                                              |
| --------- | --------------------------------------------------- |
| `↑` / `↓` | Navigate track list                                 |
| `←` / `→` | Rewind / fast-forward 5 seconds                     |
| `Enter`   | Play selected track                                 |
| `Space`   | Pause / resume                                      |
| `s`       | Stop playback                                       |
| `n` / `p` | Next / previous track                               |
| `/`       | Search                                              |
| `i`       | Track info & album art                              |
| `f`       | Fetch artwork (in info popup, when no art embedded) |
| `Esc`     | Exit search / go back / close popup                 |
| `1` - `4` | Switch view: songs, artists, albums, genres         |
| `q`       | Quit                                                |

## Features

- Recursive directory scanning for `.mp3`, `.ogg`, and `.flac` files
- Metadata display: title, artist, album, genre, duration
- Album art display with native image support (iTerm2, Ghostty, Kitty, WezTerm) and half-block fallback for other terminals
- Fetch missing album art from MusicBrainz / Cover Art Archive
- Real-time search across title, artist, and album
- Browse by artist, album, or genre with drill-down navigation
- Auto-advances to the next track when playback ends
- Single binary with no runtime configuration

## Development

```bash
make test    # Run tests
make vet     # Run go vet
make lint    # Run golangci-lint
make fmt     # Format code
make clean   # Remove binary
```

## License

MIT
