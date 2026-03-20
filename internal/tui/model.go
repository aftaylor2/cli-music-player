package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/ataylor/cli-music-player/internal/audio"
	"github.com/ataylor/cli-music-player/internal/library"
)

var highlightStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Bold(true)

// ViewMode represents the current browsing mode.
type ViewMode int

const (
	ViewSongs ViewMode = iota
	ViewArtists
	ViewAlbums
	ViewGenres
)

// tickMsg is sent periodically to update the elapsed time display.
type tickMsg time.Time

// trackDoneMsg is sent when the current track finishes.
type trackDoneMsg struct{}

// Model is the top-level Bubble Tea model for the music player.
type Model struct {
	lib    *library.Library
	player *audio.Player
	table  table.Model
	width  int
	height int

	// Search state
	searching   bool
	searchInput textinput.Model
	searchQuery string

	// View mode state
	viewMode   ViewMode
	drillGroup string // non-empty when drilled into a group
	groupKeys  []string

	// Visible tracks (filtered view)
	visibleTracks []library.Track
	trackIndex    int // index within visibleTracks of currently playing track (-1 if none)

	// Info popup
	showInfo     bool
	fetchingArt  bool
	fetchArtErr  string
}

// NewModel creates a new TUI model.
func NewModel(lib *library.Library, player *audio.Player) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100

	tracks := lib.Tracks()

	m := Model{
		lib:           lib,
		player:        player,
		table:         newTrackTable(80),
		searchInput:   ti,
		visibleTracks: tracks,
		trackIndex:    -1,
		viewMode:      ViewSongs,
	}

	m.table = m.table.WithRows(trackRows(tracks))
	return m
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table = m.table.WithTargetWidth(msg.Width).WithPageSize(msg.Height - 6)
		return m, nil

	case tickMsg:
		// Check if current track ended — auto-advance.
		if m.player.Status() == audio.Stopped && m.trackIndex >= 0 {
			next := m.trackIndex + 1
			if next < len(m.visibleTracks) {
				m.trackIndex = next
				track := m.visibleTracks[next]
				_ = m.player.Play(&track)
			} else {
				m.trackIndex = -1
			}
		}
		return m, tickCmd()

	case trackDoneMsg:
		return m, nil

	case artFetchResult:
		m.fetchingArt = false
		if msg.err != nil {
			m.fetchArtErr = msg.err.Error()
		} else if track := m.getInfoTrack(); track != nil {
			track.AlbumArt = msg.data
			track.AlbumArtMIME = "image/jpeg"
			m.fetchArtErr = ""
		}
		return m, nil

	case tea.KeyMsg:
		// If info popup is open, only i/Esc dismiss it.
		if m.showInfo {
			return m.handleInfoKey(msg)
		}
		// If searching, handle search input.
		if m.searching {
			return m.handleSearchKey(msg)
		}
		return m.handleNormalKey(msg)
	}

	// Forward to table.
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		m.searching = false
		m.searchQuery = ""
		m.searchInput.Reset()
		m.refreshVisibleTracks()
		return m, nil
	case msg.Type == tea.KeyEnter:
		m.searching = false
		return m, nil
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.searchQuery = m.searchInput.Value()
		m.refreshVisibleTracks()
		return m, cmd
	}
}

func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		m.player.Stop()
		return m, tea.Quit

	case key.Matches(msg, keys.Play):
		return m.playSelected()

	case key.Matches(msg, keys.Pause):
		m.player.TogglePause()
		return m, nil

	case key.Matches(msg, keys.Stop):
		m.player.Stop()
		m.trackIndex = -1
		return m, nil

	case key.Matches(msg, keys.Next):
		return m.playNext()

	case key.Matches(msg, keys.Prev):
		return m.playPrev()

	case key.Matches(msg, keys.SeekBack):
		m.player.Seek(-5 * time.Second)
		return m, nil

	case key.Matches(msg, keys.SeekFwd):
		m.player.Seek(5 * time.Second)
		return m, nil

	case key.Matches(msg, keys.Search):
		m.searching = true
		m.searchInput.Focus()
		return m, m.searchInput.Cursor.BlinkCmd()

	case key.Matches(msg, keys.ViewSong):
		m.viewMode = ViewSongs
		m.drillGroup = ""
		m.refreshVisibleTracks()
		return m, nil

	case key.Matches(msg, keys.ViewArtist):
		m.viewMode = ViewArtists
		m.drillGroup = ""
		m.refreshVisibleTracks()
		return m, nil

	case key.Matches(msg, keys.ViewAlbum):
		m.viewMode = ViewAlbums
		m.drillGroup = ""
		m.refreshVisibleTracks()
		return m, nil

	case key.Matches(msg, keys.ViewGenre):
		m.viewMode = ViewGenres
		m.drillGroup = ""
		m.refreshVisibleTracks()
		return m, nil

	case key.Matches(msg, keys.Info):
		track := m.getInfoTrack()
		if track != nil {
			m.showInfo = true
		}
		return m, nil

	case key.Matches(msg, keys.Escape):
		if m.drillGroup != "" {
			m.drillGroup = ""
			m.refreshVisibleTracks()
			return m, nil
		}
		return m, nil

	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
}

func (m Model) handleInfoKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Info), key.Matches(msg, keys.Escape):
		m.showInfo = false
		m.fetchingArt = false
		m.fetchArtErr = ""
		return m, nil

	case key.Matches(msg, keys.FetchArt):
		track := m.getInfoTrack()
		if track != nil && len(track.AlbumArt) == 0 && !m.fetchingArt {
			m.fetchingArt = true
			m.fetchArtErr = ""
			artist, album := track.Artist, track.Album
			return m, func() tea.Msg {
				data, err := fetchAlbumArt(artist, album)
				return artFetchResult{data: data, err: err}
			}
		}
	}
	// Swallow all other keys while popup is open.
	return m, nil
}

// getInfoTrack returns the track to show in the info popup: the currently
// playing track if available, otherwise the highlighted table row track.
func (m Model) getInfoTrack() *library.Track {
	if m.trackIndex >= 0 && m.trackIndex < len(m.visibleTracks) {
		return &m.visibleTracks[m.trackIndex]
	}
	return m.getHighlightedTrack()
}

func (m Model) getHighlightedTrack() *library.Track {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return nil
	}
	// In grouped view without drill-down, there's no individual track.
	if m.viewMode != ViewSongs && m.drillGroup == "" {
		return nil
	}
	title, _ := row.Data[colKeyTitle].(string)
	for i, t := range m.visibleTracks {
		if t.Title == title {
			return &m.visibleTracks[i]
		}
	}
	return nil
}

func (m *Model) playSelected() (tea.Model, tea.Cmd) {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return m, nil
	}

	// If in grouped view (not drilled down), drill into the group.
	if m.viewMode != ViewSongs && m.drillGroup == "" {
		title, _ := row.Data[colKeyTitle].(string)
		if title != "" {
			m.drillGroup = title
			m.refreshVisibleTracks()
		}
		return m, nil
	}

	// Find the selected track index.
	title, _ := row.Data[colKeyTitle].(string)
	for i, t := range m.visibleTracks {
		if t.Title == title {
			m.trackIndex = i
			track := m.visibleTracks[i]
			_ = m.player.Play(&track)
			break
		}
	}
	return m, nil
}

func (m *Model) playNext() (tea.Model, tea.Cmd) {
	if m.trackIndex < 0 || len(m.visibleTracks) == 0 {
		return m, nil
	}
	next := m.trackIndex + 1
	if next >= len(m.visibleTracks) {
		return m, nil
	}
	m.trackIndex = next
	track := m.visibleTracks[next]
	_ = m.player.Play(&track)
	return m, nil
}

func (m *Model) playPrev() (tea.Model, tea.Cmd) {
	if m.trackIndex <= 0 || len(m.visibleTracks) == 0 {
		return m, nil
	}
	prev := m.trackIndex - 1
	m.trackIndex = prev
	track := m.visibleTracks[prev]
	_ = m.player.Play(&track)
	return m, nil
}

func (m *Model) refreshVisibleTracks() {
	var tracks []library.Track

	switch m.viewMode {
	case ViewSongs:
		tracks = m.lib.Tracks()
	case ViewArtists:
		if m.drillGroup != "" {
			tracks = m.lib.FilterByArtist(m.drillGroup)
		}
	case ViewAlbums:
		if m.drillGroup != "" {
			tracks = m.lib.FilterByAlbum(m.drillGroup)
		}
	case ViewGenres:
		if m.drillGroup != "" {
			tracks = filterByGenre(m.lib, m.drillGroup)
		}
	}

	// Apply search filter.
	if m.searchQuery != "" && len(tracks) > 0 {
		tracks = searchTracks(tracks, m.searchQuery)
	}

	// If in grouped view and not drilled down, show group rows.
	if m.viewMode != ViewSongs && m.drillGroup == "" {
		m.showGroupView()
		return
	}

	if m.viewMode != ViewSongs && m.drillGroup == "" {
		return
	}

	// For song view with search.
	if m.viewMode == ViewSongs && m.searchQuery != "" {
		tracks = m.lib.Search(m.searchQuery)
	} else if m.viewMode == ViewSongs {
		tracks = m.lib.Tracks()
	}

	m.visibleTracks = tracks
	m.table = m.table.WithRows(trackRows(tracks))
}

func (m *Model) showGroupView() {
	var groups map[string][]library.Track
	switch m.viewMode {
	case ViewArtists:
		groups = m.lib.GroupByArtist()
	case ViewAlbums:
		groups = m.lib.GroupByAlbum()
	case ViewGenres:
		groups = m.lib.GroupByGenre()
	default:
		return
	}

	m.groupKeys = sortedKeys(groups)
	m.table = m.table.WithRows(groupRows(groups, m.groupKeys))
}

func filterByGenre(lib *library.Library, genre string) []library.Track {
	var result []library.Track
	for _, t := range lib.Tracks() {
		if t.Genre == genre {
			result = append(result, t)
		}
	}
	return result
}

func searchTracks(tracks []library.Track, query string) []library.Track {
	lib := library.NewLibraryFromTracks("", tracks)
	return lib.Search(query)
}

func sortedKeys(m map[string][]library.Track) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort alphabetically.
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func (m Model) View() string {
	if m.showInfo {
		track := m.getInfoTrack()
		return renderInfoPopup(track, m.width, m.height, m.fetchingArt, m.fetchArtErr)
	}

	var header string
	switch m.viewMode {
	case ViewSongs:
		header = "Songs"
	case ViewArtists:
		if m.drillGroup != "" {
			header = "Artist: " + m.drillGroup
		} else {
			header = "Artists"
		}
	case ViewAlbums:
		if m.drillGroup != "" {
			header = "Album: " + m.drillGroup
		} else {
			header = "Albums"
		}
	case ViewGenres:
		if m.drillGroup != "" {
			header = "Genre: " + m.drillGroup
		} else {
			header = "Genres"
		}
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	view := headerStyle.Render(header) + "\n"

	if m.searching {
		view += m.searchInput.View() + "\n"
	}

	view += m.table.View() + "\n"
	view += renderControls(
		m.player.CurrentTrack(),
		m.player.Status(),
		m.player.Elapsed(),
		m.width,
	)
	view += "\n"

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	view += helpStyle.Render("↑↓ navigate  ←→ seek 5s  enter play  space pause  s stop  n/p next/prev  i info  / search  1-4 views  q quit")

	return view
}
