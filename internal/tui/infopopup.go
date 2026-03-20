package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ataylor/cli-music-player/internal/library"
)

var (
	popupBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212")).
				Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

var (
	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	fetchingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
)

func renderInfoPopup(track *library.Track, width, height int, fetchingArt bool, fetchErr string) string {
	if track == nil {
		return ""
	}

	// Reserve space for border + padding.
	innerWidth := width - 8
	innerHeight := height - 6
	if innerWidth < 20 {
		innerWidth = 20
	}
	if innerHeight < 5 {
		innerHeight = 5
	}

	// Determine available width for metadata values.
	// In two-column mode, metadata gets roughly half the inner width minus the gap.
	// In single-column mode, it gets the full inner width.
	metaMaxWidth := innerWidth
	showArt := width >= 50 && height >= 15
	if showArt {
		artWidth := 36
		if artWidth > innerWidth/2 {
			artWidth = innerWidth / 2
		}
		metaMaxWidth = innerWidth - artWidth - 3
		if metaMaxWidth < 15 {
			metaMaxWidth = 15
		}
	}
	// Reserve space for the longest label ("Duration:") + space = ~10 chars.
	maxValueLen := metaMaxWidth - 12
	if maxValueLen < 10 {
		maxValueLen = 10
	}

	// Build metadata lines.
	meta := []struct{ label, value string }{
		{"Title", track.Title},
		{"Artist", track.Artist},
		{"Album", track.Album},
		{"Genre", track.Genre},
		{"Track #", fmt.Sprintf("%d", track.TrackNumber)},
		{"Duration", track.FormatDuration()},
		{"Format", track.Format.String()},
		{"File", truncate(track.FilePath, maxValueLen)},
	}

	var metaLines []string
	for _, m := range meta {
		line := labelStyle.Render(m.label+":") + " " + valueStyle.Render(m.value)
		metaLines = append(metaLines, line)
	}

	// Show fetch status/hint when no embedded art.
	if len(track.AlbumArt) == 0 {
		if fetchingArt {
			metaLines = append(metaLines, "", fetchingStyle.Render("Fetching artwork..."))
		} else if fetchErr != "" {
			metaLines = append(metaLines, "", errorStyle.Render("Fetch failed: "+fetchErr))
			metaLines = append(metaLines, hintStyle.Render("Press f to retry"))
		} else {
			metaLines = append(metaLines, "", hintStyle.Render("Press f to fetch artwork"))
		}
	}

	metaBlock := strings.Join(metaLines, "\n")

	var content string

	if showArt {
		artWidth := 36
		if artWidth > innerWidth/2 {
			artWidth = innerWidth / 2
		}
		art := renderAlbumArt(track.AlbumArt, artWidth)

		artBox := lipgloss.NewStyle().Width(artWidth).Render(art)
		metaBox := lipgloss.NewStyle().Width(metaMaxWidth).MaxWidth(metaMaxWidth).Render(metaBlock)
		content = lipgloss.JoinHorizontal(lipgloss.Top, artBox, "   ", metaBox)
	} else {
		// Narrow terminal: metadata only.
		content = metaBlock
	}

	popup := popupBorderStyle.
		Width(innerWidth).
		MaxHeight(innerHeight).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}
