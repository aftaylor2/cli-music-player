package tui

import (
	"fmt"

	"github.com/evertras/bubble-table/table"

	"github.com/ataylor/cli-music-player/internal/library"
)

const (
	colKeyNum      = "num"
	colKeyTitle    = "title"
	colKeyArtist   = "artist"
	colKeyAlbum    = "album"
	colKeyGenre    = "genre"
	colKeyDuration = "duration"
)

func newTrackTable(width int) table.Model {
	columns := []table.Column{
		table.NewColumn(colKeyNum, "#", 5),
		table.NewColumn(colKeyTitle, "Title", 30),
		table.NewColumn(colKeyArtist, "Artist", 22),
		table.NewColumn(colKeyAlbum, "Album", 22),
		table.NewColumn(colKeyGenre, "Genre", 12),
		table.NewColumn(colKeyDuration, "Duration", 9),
	}

	return table.New(columns).
		WithPageSize(20).
		WithTargetWidth(width).
		Focused(true).
		HighlightStyle(highlightStyle)
}

func trackRows(tracks []library.Track) []table.Row {
	rows := make([]table.Row, len(tracks))
	for i, t := range tracks {
		rows[i] = table.NewRow(table.RowData{
			colKeyNum:      fmt.Sprintf("%d", i+1),
			colKeyTitle:    t.Title,
			colKeyArtist:   t.Artist,
			colKeyAlbum:    t.Album,
			colKeyGenre:    t.Genre,
			colKeyDuration: t.FormatDuration(),
		})
	}
	return rows
}

func groupRows(groups map[string][]library.Track, keys []string) []table.Row {
	rows := make([]table.Row, len(keys))
	for i, k := range keys {
		rows[i] = table.NewRow(table.RowData{
			colKeyNum:      fmt.Sprintf("%d", i+1),
			colKeyTitle:    k,
			colKeyArtist:   "",
			colKeyAlbum:    "",
			colKeyGenre:    "",
			colKeyDuration: fmt.Sprintf("%d tracks", len(groups[k])),
		})
	}
	return rows
}
