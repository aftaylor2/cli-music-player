package library

import (
	"testing"
	"time"
)

func makeTracks() []Track {
	return []Track{
		{FilePath: "/music/a.mp3", Title: "Alpha", Artist: "ArtistA", Album: "Album1", Genre: "Rock", Duration: 3 * time.Minute},
		{FilePath: "/music/b.ogg", Title: "Beta", Artist: "ArtistB", Album: "Album2", Genre: "Jazz", Duration: 4 * time.Minute},
		{FilePath: "/music/c.mp3", Title: "Charlie", Artist: "ArtistA", Album: "Album1", Genre: "Rock", Duration: 5 * time.Minute},
		{FilePath: "/music/d.ogg", Title: "Delta", Artist: "ArtistC", Album: "Album3", Genre: "Jazz", Duration: 2 * time.Minute},
	}
}

func TestLibraryFromTracks(t *testing.T) {
	tracks := makeTracks()
	lib := NewLibraryFromTracks("/music", tracks)

	if lib.Len() != 4 {
		t.Errorf("Len() = %d, want 4", lib.Len())
	}
	if lib.RootDir != "/music" {
		t.Errorf("RootDir = %q, want %q", lib.RootDir, "/music")
	}
}

func TestLibraryEmpty(t *testing.T) {
	lib := NewLibraryFromTracks("/empty", nil)
	if lib.Len() != 0 {
		t.Errorf("Len() = %d, want 0", lib.Len())
	}
	if len(lib.Tracks()) != 0 {
		t.Errorf("Tracks() returned %d items, want 0", len(lib.Tracks()))
	}
}

func TestLibraryTracks(t *testing.T) {
	tracks := makeTracks()
	lib := NewLibraryFromTracks("/music", tracks)

	got := lib.Tracks()
	if len(got) != len(tracks) {
		t.Errorf("Tracks() returned %d items, want %d", len(got), len(tracks))
	}
}

// T020: Search tests
func TestSearchEmptyQuery(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("")
	if len(got) != 4 {
		t.Errorf("Search(\"\") = %d results, want 4", len(got))
	}
}

func TestSearchByTitle(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("alpha")
	if len(got) != 1 {
		t.Errorf("Search(\"alpha\") = %d results, want 1", len(got))
	}
	if len(got) > 0 && got[0].Title != "Alpha" {
		t.Errorf("Search(\"alpha\")[0].Title = %q, want %q", got[0].Title, "Alpha")
	}
}

func TestSearchByArtist(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("ArtistA")
	if len(got) != 2 {
		t.Errorf("Search(\"ArtistA\") = %d results, want 2", len(got))
	}
}

func TestSearchByAlbum(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("Album2")
	if len(got) != 1 {
		t.Errorf("Search(\"Album2\") = %d results, want 1", len(got))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("CHARLIE")
	if len(got) != 1 {
		t.Errorf("Search(\"CHARLIE\") = %d results, want 1", len(got))
	}
}

func TestSearchNoResults(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.Search("zzzzz")
	if len(got) != 0 {
		t.Errorf("Search(\"zzzzz\") = %d results, want 0", len(got))
	}
}

// T024: GroupBy and FilterBy tests
func TestGroupByArtist(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	groups := lib.GroupByArtist()
	if len(groups) != 3 {
		t.Errorf("GroupByArtist() = %d groups, want 3", len(groups))
	}
	if len(groups["ArtistA"]) != 2 {
		t.Errorf("GroupByArtist()[\"ArtistA\"] = %d tracks, want 2", len(groups["ArtistA"]))
	}
}

func TestGroupByAlbum(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	groups := lib.GroupByAlbum()
	if len(groups) != 3 {
		t.Errorf("GroupByAlbum() = %d groups, want 3", len(groups))
	}
}

func TestGroupByGenre(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	groups := lib.GroupByGenre()
	if len(groups) != 2 {
		t.Errorf("GroupByGenre() = %d groups, want 2", len(groups))
	}
	if len(groups["Rock"]) != 2 {
		t.Errorf("GroupByGenre()[\"Rock\"] = %d tracks, want 2", len(groups["Rock"]))
	}
}

func TestFilterByArtist(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.FilterByArtist("ArtistA")
	if len(got) != 2 {
		t.Errorf("FilterByArtist(\"ArtistA\") = %d results, want 2", len(got))
	}
}

func TestFilterByAlbum(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.FilterByAlbum("Album1")
	if len(got) != 2 {
		t.Errorf("FilterByAlbum(\"Album1\") = %d results, want 2", len(got))
	}
}

func TestFilterByArtistNoMatch(t *testing.T) {
	lib := NewLibraryFromTracks("/music", makeTracks())
	got := lib.FilterByArtist("Nobody")
	if len(got) != 0 {
		t.Errorf("FilterByArtist(\"Nobody\") = %d results, want 0", len(got))
	}
}
