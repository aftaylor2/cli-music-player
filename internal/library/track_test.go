package library

import (
	"testing"
	"time"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    AudioFormat
		wantErr bool
	}{
		{"mp3 lowercase", "song.mp3", FormatMP3, false},
		{"mp3 uppercase", "song.MP3", FormatMP3, false},
		{"mp3 mixed case", "song.Mp3", FormatMP3, false},
		{"ogg lowercase", "song.ogg", FormatOGG, false},
		{"ogg uppercase", "song.OGG", FormatOGG, false},
		{"unsupported wav", "song.wav", 0, true},
		{"flac lowercase", "song.flac", FormatFLAC, false},
		{"flac uppercase", "song.FLAC", FormatFLAC, false},
		{"no extension", "song", 0, true},
		{"path with dirs", "/home/user/music/song.mp3", FormatMP3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectFormat(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFormat(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("DetectFormat(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"zero", 0, "0:00"},
		{"thirty seconds", 30 * time.Second, "0:30"},
		{"one minute", 60 * time.Second, "1:00"},
		{"three minutes ten", 3*time.Minute + 10*time.Second, "3:10"},
		{"one hour", time.Hour, "1:00:00"},
		{"one hour five min", time.Hour + 5*time.Minute + 30*time.Second, "1:05:30"},
		{"negative treated as zero", -5 * time.Second, "0:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Track{Duration: tt.duration}
			got := tr.FormatDuration()
			if got != tt.want {
				t.Errorf("FormatDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFilenameWithoutExt(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"song.mp3", "song"},
		{"/path/to/My Song.ogg", "My Song"},
		{"no-ext", "no-ext"},
		{"/dir/file.tar.gz", "file.tar"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := filenameWithoutExt(tt.path)
			if got != tt.want {
				t.Errorf("filenameWithoutExt(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestAudioFormatString(t *testing.T) {
	if FormatMP3.String() != "MP3" {
		t.Errorf("FormatMP3.String() = %q, want %q", FormatMP3.String(), "MP3")
	}
	if FormatOGG.String() != "OGG" {
		t.Errorf("FormatOGG.String() = %q, want %q", FormatOGG.String(), "OGG")
	}
	if FormatFLAC.String() != "FLAC" {
		t.Errorf("FormatFLAC.String() = %q, want %q", FormatFLAC.String(), "FLAC")
	}
}
