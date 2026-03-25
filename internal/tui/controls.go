package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/aftaylor2/cli-music-player/internal/audio"
	"github.com/aftaylor2/cli-music-player/internal/library"
)

var (
	controlsStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("240"))

	trackTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	trackArtistStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("251"))

	statusIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
)

func renderControls(track *library.Track, status audio.PlayStatus, elapsed time.Duration, width int) string {
	if track == nil {
		return controlsStyle.Width(width).Render("No track playing")
	}

	icon := statusIcon(status)
	title := trackTitleStyle.Render(truncate(track.Title, 40))
	artist := trackArtistStyle.Render(truncate(track.Artist, 30))
	elapsedStr := formatDuration(elapsed)
	totalStr := formatDuration(track.Duration)
	timeInfo := timeStyle.Render(fmt.Sprintf("%s / %s", elapsedStr, totalStr))

	line := fmt.Sprintf("%s  %s - %s    %s", icon, title, artist, timeInfo)
	return controlsStyle.Width(width).Render(line)
}

func statusIcon(s audio.PlayStatus) string {
	switch s {
	case audio.Playing:
		return statusIconStyle.Render("▶")
	case audio.Paused:
		return statusIconStyle.Render("⏸")
	default:
		return statusIconStyle.Render("⏹")
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
