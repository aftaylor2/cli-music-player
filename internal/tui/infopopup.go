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

	innerWidth := width - 8
	innerHeight := height - 6
	if innerWidth < 20 {
		innerWidth = 20
	}
	if innerHeight < 5 {
		innerHeight = 5
	}

	showArt := width >= 50 && height >= 15
	proto := detectImageProtocol()
	useNativeImage := showArt && len(track.AlbumArt) > 0 && proto != protoHalfBlock

	// For native image protocols, render the popup with metadata only,
	// then use cursor movement to overlay the image at the right position.
	if useNativeImage {
		return renderNativeInfoPopup(track, width, height, innerWidth, innerHeight, fetchingArt, fetchErr)
	}

	useHalfBlock := showArt && !useNativeImage

	metaMaxWidth := innerWidth
	if useHalfBlock {
		artWidth := artColumns(innerWidth)
		metaMaxWidth = innerWidth - artWidth - 3
		if metaMaxWidth < 15 {
			metaMaxWidth = 15
		}
	}

	metaBlock := buildMetaBlock(track, metaMaxWidth, fetchingArt, fetchErr)

	var content string
	if useHalfBlock {
		artWidth := artColumns(innerWidth)
		art := renderAlbumArt(track.AlbumArt, artWidth)
		artBox := lipgloss.NewStyle().Width(artWidth).Render(art)
		metaBox := lipgloss.NewStyle().Width(metaMaxWidth).MaxWidth(metaMaxWidth).Render(metaBlock)
		content = lipgloss.JoinHorizontal(lipgloss.Top, artBox, "   ", metaBox)
	} else {
		content = metaBlock
	}

	popup := popupBorderStyle.
		Width(innerWidth).
		MaxHeight(innerHeight).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}

// renderNativeInfoPopup builds the info popup for iTerm2/Kitty by rendering
// the popup with metadata only via lipgloss, then using ANSI cursor movement
// to overlay the image inside the popup area.
func renderNativeInfoPopup(track *library.Track, width, height, innerWidth, innerHeight int, fetchingArt bool, fetchErr string) string {
	artWidth := artColumns(innerWidth)
	metaColWidth := innerWidth - artWidth - 3
	if metaColWidth < 15 {
		metaColWidth = 15
	}

	metaBlock := buildMetaBlock(track, metaColWidth, fetchingArt, fetchErr)

	// Compute image row count for the placeholder height.
	img, err := decodeAlbumArt(track.AlbumArt)
	artRows := artWidth / 2
	if err == nil {
		b := img.Bounds()
		artRows = imageRowCount(b.Dx(), b.Dy(), artWidth)
	}

	// Build a blank placeholder the same size as the image.
	spaceLine := strings.Repeat(" ", artWidth)
	var placeholderLines []string
	for i := 0; i < artRows; i++ {
		placeholderLines = append(placeholderLines, spaceLine)
	}
	placeholder := strings.Join(placeholderLines, "\n")

	// Use JoinHorizontal — same as half-block path but left column is blank.
	artBox := lipgloss.NewStyle().Width(artWidth).Render(placeholder)
	metaBox := lipgloss.NewStyle().Width(metaColWidth).MaxWidth(metaColWidth).Render(metaBlock)
	content := lipgloss.JoinHorizontal(lipgloss.Top, artBox, "   ", metaBox)

	popup := popupBorderStyle.
		Width(innerWidth).
		MaxHeight(innerHeight).
		Render(content)

	placed := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)

	// Calculate image position mathematically from known dimensions.
	// lipgloss Width() sets content+padding width (excludes border).
	// Border adds 1 char on each side.
	popupVisualWidth := innerWidth + 2
	leftMargin := (width - popupVisualWidth) / 2

	// Content starts after border(1) + padding(2) from left edge.
	imgCol := leftMargin + 1 + 2

	// Find the popup top row by searching for the border character.
	lines := strings.Split(placed, "\n")
	popupTopRow := 0
	for i, line := range lines {
		if strings.Contains(line, "╭") {
			popupTopRow = i
			break
		}
	}
	// Content starts after border(1) + padding(1) rows.
	imgRow := popupTopRow + 2

	artEsc := renderAlbumArt(track.AlbumArt, artWidth)

	// ANSI cursor save/move/restore to overlay the image.
	overlay := fmt.Sprintf("\x1b[s\x1b[%d;%dH%s\x1b[u", imgRow+1, imgCol+1, artEsc)

	return placed + overlay
}

func buildMetaBlock(track *library.Track, maxWidth int, fetchingArt bool, fetchErr string) string {
	maxValueLen := maxWidth - 12
	if maxValueLen < 10 {
		maxValueLen = 10
	}

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

	return strings.Join(metaLines, "\n")
}

func artColumns(innerWidth int) int {
	w := 36
	if w > innerWidth/2 {
		w = innerWidth / 2
	}
	return w
}
