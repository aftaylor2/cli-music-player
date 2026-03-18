# CLI Interface Contract

## Command

```
player [directory]
```

## Arguments

| Argument    | Required | Default | Description                     |
|-------------|----------|---------|---------------------------------|
| `directory` | No       | `.`     | Path to directory to scan for audio files |

## Exit Codes

| Code | Meaning                                    |
|------|--------------------------------------------|
| 0    | Normal exit (user quit)                    |
| 1    | Error: specified directory does not exist   |
| 2    | Error: specified path is not a directory    |

## Supported File Types

Files are matched by extension (case-insensitive):

| Extension | Format     |
|-----------|------------|
| `.mp3`    | MPEG Audio |
| `.ogg`    | Ogg Vorbis |

## Behavior

1. On launch, recursively scan `directory` for files matching
   supported extensions.
2. Parse metadata tags from each file.
3. Render TUI with track table and playback controls.
4. On exit (user presses `q`), stop playback and restore terminal.

## Error Output

Errors are written to stderr. The TUI occupies stdout via the
alternate screen buffer (Bubble Tea default behavior).
