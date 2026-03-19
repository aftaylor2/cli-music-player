package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
)

// AudioFormat represents a supported audio file format.
type AudioFormat int

const (
	FormatMP3 AudioFormat = iota
	FormatOGG
	FormatFLAC
)

func (f AudioFormat) String() string {
	switch f {
	case FormatMP3:
		return "MP3"
	case FormatOGG:
		return "OGG"
	case FormatFLAC:
		return "FLAC"
	default:
		return "Unknown"
	}
}

// Track represents a single audio file with its metadata.
type Track struct {
	FilePath    string
	Title       string
	Artist      string
	Album       string
	Genre       string
	TrackNumber int
	Duration    time.Duration
	Format      AudioFormat
}

// FormatDuration returns the track duration as "m:ss" or "h:mm:ss".
func (t Track) FormatDuration() string {
	d := t.Duration
	if d < 0 {
		d = 0
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

// DetectFormat returns the AudioFormat for a file path based on its extension,
// or an error if the format is unsupported.
func DetectFormat(path string) (AudioFormat, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp3":
		return FormatMP3, nil
	case ".ogg":
		return FormatOGG, nil
	case ".flac":
		return FormatFLAC, nil
	default:
		return 0, fmt.Errorf("unsupported format: %s", ext)
	}
}

// NewTrackFromFile reads metadata and duration from an audio file and returns
// a Track. Missing metadata fields fall back to sensible defaults.
func NewTrackFromFile(path string) (Track, error) {
	format, err := DetectFormat(path)
	if err != nil {
		return Track{}, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return Track{}, fmt.Errorf("resolving path: %w", err)
	}

	t := Track{
		FilePath: absPath,
		Format:   format,
		Title:    filenameWithoutExt(path),
		Artist:   "Unknown",
		Album:    "Unknown",
		Genre:    "Unknown",
	}

	f, err := os.Open(absPath)
	if err != nil {
		return Track{}, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	// Read metadata tags.
	m, err := tag.ReadFrom(f)
	if err == nil {
		if v := strings.TrimSpace(m.Title()); v != "" {
			t.Title = v
		}
		if v := strings.TrimSpace(m.Artist()); v != "" {
			t.Artist = v
		}
		if v := strings.TrimSpace(m.Album()); v != "" {
			t.Album = v
		}
		if v := strings.TrimSpace(m.Genre()); v != "" {
			t.Genre = v
		}
		num, _ := m.Track()
		t.TrackNumber = num
	}
	// If tag reading fails, we keep fallback values.

	// Compute duration by decoding the stream header.
	t.Duration = readDuration(absPath, format)

	return t, nil
}

func readDuration(path string, format AudioFormat) time.Duration {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var af beep.Format

	switch format {
	case FormatMP3:
		streamer, af, err = mp3.Decode(f)
	case FormatOGG:
		streamer, af, err = vorbis.Decode(f)
	case FormatFLAC:
		streamer, af, err = flac.Decode(f)
	}
	if err != nil || streamer == nil {
		return 0
	}
	defer streamer.Close()

	totalSamples := streamer.Len()
	if totalSamples <= 0 || af.SampleRate <= 0 {
		return 0
	}
	return af.SampleRate.D(totalSamples)
}

func filenameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
