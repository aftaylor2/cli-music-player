package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Play       key.Binding
	Pause      key.Binding
	Stop       key.Binding
	Next       key.Binding
	Prev       key.Binding
	Search     key.Binding
	Escape     key.Binding
	Quit       key.Binding
	ViewSong   key.Binding
	ViewArtist key.Binding
	ViewAlbum  key.Binding
	ViewGenre  key.Binding
	SeekBack   key.Binding
	SeekFwd    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	Play: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause/resume"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back/cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	ViewSong: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "songs"),
	),
	ViewArtist: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "artists"),
	),
	ViewAlbum: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "albums"),
	),
	ViewGenre: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "genres"),
	),
	SeekBack: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "rewind 5s"),
	),
	SeekFwd: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "forward 5s"),
	),
}
