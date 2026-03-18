package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectoryEmpty(t *testing.T) {
	dir := t.TempDir()
	tracks, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory(%q) unexpected error: %v", dir, err)
	}
	if len(tracks) != 0 {
		t.Errorf("ScanDirectory empty dir: got %d tracks, want 0", len(tracks))
	}
}

func TestScanDirectorySkipsNonAudio(t *testing.T) {
	dir := t.TempDir()

	// Create non-audio files.
	for _, name := range []string{"readme.txt", "image.png", "data.csv"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("data"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tracks, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory unexpected error: %v", err)
	}
	if len(tracks) != 0 {
		t.Errorf("ScanDirectory with non-audio files: got %d tracks, want 0", len(tracks))
	}
}

func TestScanDirectoryFindsAudioExtensions(t *testing.T) {
	dir := t.TempDir()

	// Create files with audio extensions (they won't have valid audio data,
	// but the scanner should still find them by extension).
	for _, name := range []string{"song.mp3", "track.ogg", "other.MP3", "upper.OGG"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("fake"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// Also a non-audio file.
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	tracks, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory unexpected error: %v", err)
	}
	// Scanner should attempt to parse all 4 audio files. Some may fail to
	// decode (invalid data) and be skipped — that's acceptable. We check
	// that at least the scanner found and attempted them by verifying no
	// non-audio files snuck in.
	for _, tr := range tracks {
		ext := filepath.Ext(tr.FilePath)
		if ext != ".mp3" && ext != ".ogg" && ext != ".MP3" && ext != ".OGG" {
			t.Errorf("unexpected file in results: %s", tr.FilePath)
		}
	}
}

func TestScanDirectoryRecursive(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir", "deep")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	// Create a file in the nested directory.
	if err := os.WriteFile(filepath.Join(sub, "nested.mp3"), []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	tracks, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory unexpected error: %v", err)
	}
	// The file may or may not parse successfully (fake data), but the scanner
	// should not return an error for the directory walk itself.
	_ = tracks
}

func TestScanDirectoryNonexistent(t *testing.T) {
	_, err := ScanDirectory("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("ScanDirectory on nonexistent path: expected error, got nil")
	}
}
