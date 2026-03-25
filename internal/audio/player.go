package audio

import (
	"fmt"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"

	"github.com/aftaylor2/cli-music-player/internal/library"
)

// PlayStatus represents the current state of the audio player.
type PlayStatus int

const (
	Stopped PlayStatus = iota
	Playing
	Paused
)

func (s PlayStatus) String() string {
	switch s {
	case Stopped:
		return "Stopped"
	case Playing:
		return "Playing"
	case Paused:
		return "Paused"
	default:
		return "Unknown"
	}
}

// Player manages audio playback for a single track at a time.
type Player struct {
	mu           sync.Mutex
	status       PlayStatus
	currentTrack *library.Track
	ctrl         *beep.Ctrl
	streamer     beep.StreamSeekCloser
	sampleRate   beep.SampleRate
	speakerInit  bool
	done         chan struct{}
}

// New creates a new Player instance.
func New() *Player {
	return &Player{
		status: Stopped,
	}
}

// Play starts playback of the given track. If another track is currently
// playing, it is stopped first.
func (p *Player) Play(track *library.Track) error {
	p.Stop()

	p.mu.Lock()
	defer p.mu.Unlock()

	streamer, af, err := OpenDecoder(track.FilePath, track.Format)
	if err != nil {
		return err
	}

	sr := af.SampleRate

	if !p.speakerInit {
		err = speaker.Init(sr, sr.N(time.Millisecond*100))
		if err != nil {
			streamer.Close()
			return fmt.Errorf("initializing speaker: %w", err)
		}
		p.speakerInit = true
		p.sampleRate = sr
	}

	// Resample if the track's sample rate differs from the speaker's.
	var playStream beep.Streamer
	if sr != p.sampleRate {
		playStream = beep.Resample(4, sr, p.sampleRate, streamer)
	} else {
		playStream = streamer
	}

	p.streamer = streamer
	p.ctrl = &beep.Ctrl{Streamer: playStream, Paused: false}
	p.currentTrack = track
	p.status = Playing
	p.done = make(chan struct{})

	done := p.done
	speaker.Play(beep.Seq(p.ctrl, beep.Callback(func() {
		p.mu.Lock()
		p.status = Stopped
		p.mu.Unlock()
		close(done)
	})))

	return nil
}

// Pause pauses the current playback.
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == Playing && p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
		p.status = Paused
	}
}

// Resume resumes paused playback.
func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.status == Paused && p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = false
		speaker.Unlock()
		p.status = Playing
	}
}

// TogglePause toggles between playing and paused states.
func (p *Player) TogglePause() {
	p.mu.Lock()
	status := p.status
	p.mu.Unlock()

	switch status {
	case Playing:
		p.Pause()
	case Paused:
		p.Resume()
	}
}

// Stop stops the current playback and releases resources.
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		speaker.Unlock()
	}

	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}

	p.ctrl = nil
	p.currentTrack = nil
	p.status = Stopped
}

// Status returns the current playback status.
func (p *Player) Status() PlayStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.status
}

// CurrentTrack returns the currently loaded track, or nil if stopped.
func (p *Player) CurrentTrack() *library.Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentTrack
}

// Seek adjusts the playback position by the given duration. Positive values
// seek forward, negative values seek backward. The position is clamped to
// [0, stream length].
func (p *Player) Seek(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.streamer == nil || p.sampleRate == 0 {
		return
	}
	speaker.Lock()
	pos := p.streamer.Position()
	total := p.streamer.Len()
	speaker.Unlock()

	delta := p.sampleRate.N(d)
	newPos := pos + delta
	if newPos < 0 {
		newPos = 0
	}
	if newPos > total {
		newPos = total
	}

	speaker.Lock()
	_ = p.streamer.Seek(newPos)
	speaker.Unlock()
}

// Elapsed returns the current playback position.
func (p *Player) Elapsed() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.streamer == nil || p.sampleRate == 0 {
		return 0
	}
	speaker.Lock()
	pos := p.streamer.Position()
	speaker.Unlock()
	return p.sampleRate.D(pos)
}

// Done returns a channel that is closed when the current track finishes.
// Returns nil if no track is playing.
func (p *Player) Done() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.done
}
