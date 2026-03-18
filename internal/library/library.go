package library

import (
	"sort"
	"strings"
)

// Library is the in-memory collection of all scanned tracks.
type Library struct {
	RootDir string
	tracks  []Track
}

// NewLibrary scans rootDir for audio files and returns a Library.
func NewLibrary(rootDir string) (*Library, error) {
	tracks, err := ScanDirectory(rootDir)
	if err != nil {
		return nil, err
	}
	sort.Slice(tracks, func(i, j int) bool {
		if tracks[i].Artist != tracks[j].Artist {
			return tracks[i].Artist < tracks[j].Artist
		}
		if tracks[i].Album != tracks[j].Album {
			return tracks[i].Album < tracks[j].Album
		}
		if tracks[i].TrackNumber != tracks[j].TrackNumber {
			return tracks[i].TrackNumber < tracks[j].TrackNumber
		}
		return tracks[i].Title < tracks[j].Title
	})
	return &Library{RootDir: rootDir, tracks: tracks}, nil
}

// NewLibraryFromTracks creates a Library from a pre-built track list.
func NewLibraryFromTracks(rootDir string, tracks []Track) *Library {
	return &Library{RootDir: rootDir, tracks: tracks}
}

// Tracks returns all tracks in the library.
func (l *Library) Tracks() []Track {
	return l.tracks
}

// Len returns the number of tracks.
func (l *Library) Len() int {
	return len(l.tracks)
}

// Search returns tracks where the query is a case-insensitive substring of
// title, artist, or album. An empty query returns all tracks.
func (l *Library) Search(query string) []Track {
	if query == "" {
		return l.tracks
	}
	q := strings.ToLower(query)
	var results []Track
	for _, t := range l.tracks {
		if strings.Contains(strings.ToLower(t.Title), q) ||
			strings.Contains(strings.ToLower(t.Artist), q) ||
			strings.Contains(strings.ToLower(t.Album), q) {
			results = append(results, t)
		}
	}
	return results
}

// GroupByArtist returns tracks grouped by artist name.
func (l *Library) GroupByArtist() map[string][]Track {
	return groupBy(l.tracks, func(t Track) string { return t.Artist })
}

// GroupByAlbum returns tracks grouped by album name.
func (l *Library) GroupByAlbum() map[string][]Track {
	return groupBy(l.tracks, func(t Track) string { return t.Album })
}

// GroupByGenre returns tracks grouped by genre.
func (l *Library) GroupByGenre() map[string][]Track {
	return groupBy(l.tracks, func(t Track) string { return t.Genre })
}

// FilterByArtist returns tracks matching the given artist.
func (l *Library) FilterByArtist(artist string) []Track {
	return filterBy(l.tracks, func(t Track) bool { return t.Artist == artist })
}

// FilterByAlbum returns tracks matching the given album.
func (l *Library) FilterByAlbum(album string) []Track {
	return filterBy(l.tracks, func(t Track) bool { return t.Album == album })
}

func groupBy(tracks []Track, key func(Track) string) map[string][]Track {
	m := make(map[string][]Track)
	for _, t := range tracks {
		k := key(t)
		m[k] = append(m[k], t)
	}
	return m
}

func filterBy(tracks []Track, pred func(Track) bool) []Track {
	var result []Track
	for _, t := range tracks {
		if pred(t) {
			result = append(result, t)
		}
	}
	return result
}
