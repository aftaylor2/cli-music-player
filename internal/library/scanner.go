package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScanDirectory recursively walks root and returns a Track for every supported
// audio file found. Files that cannot be parsed are silently skipped.
func ScanDirectory(root string) ([]Track, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("accessing directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	var tracks []Track

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip entries we can't access
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".mp3" && ext != ".ogg" {
			return nil
		}

		track, err := NewTrackFromFile(path)
		if err != nil {
			// Skip files that can't be parsed (corrupt, unreadable, etc.)
			return nil
		}

		tracks = append(tracks, track)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning directory: %w", err)
	}

	return tracks, nil
}
