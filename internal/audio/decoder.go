package audio

import (
	"fmt"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"

	"github.com/ataylor/cli-music-player/internal/library"
)

// OpenDecoder opens an audio file and returns a decoded stream along with its
// format. The caller is responsible for closing the returned streamer.
// The underlying file is kept open and will be closed when the streamer is closed.
func OpenDecoder(path string, format library.AudioFormat) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("opening audio file: %w", err)
	}

	var streamer beep.StreamSeekCloser
	var af beep.Format

	switch format {
	case library.FormatMP3:
		streamer, af, err = mp3.Decode(f)
	case library.FormatOGG:
		streamer, af, err = vorbis.Decode(f)
	default:
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("unsupported audio format: %s", format)
	}

	if err != nil {
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("decoding audio: %w", err)
	}

	return streamer, af, nil
}
