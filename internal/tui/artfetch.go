package tui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// artFetchResult holds the result of an album art fetch attempt.
type artFetchResult struct {
	data []byte
	err  error
}

// fetchAlbumArt searches MusicBrainz for a release matching artist+album,
// then downloads the front cover from the Cover Art Archive.
func fetchAlbumArt(artist, album string) ([]byte, error) {
	if artist == "" || album == "" || artist == "Unknown" || album == "Unknown" {
		return nil, fmt.Errorf("missing artist or album metadata")
	}

	mbid, err := searchMusicBrainz(artist, album)
	if err != nil {
		return nil, fmt.Errorf("MusicBrainz search: %w", err)
	}

	data, err := downloadCoverArt(mbid)
	if err != nil {
		return nil, fmt.Errorf("Cover Art Archive: %w", err)
	}
	return data, nil
}

func searchMusicBrainz(artist, album string) (string, error) {
	query := fmt.Sprintf("artist:%s AND release:%s",
		url.QueryEscape(artist), url.QueryEscape(album))
	reqURL := fmt.Sprintf("https://musicbrainz.org/ws/2/release/?query=%s&limit=1&fmt=json", url.QueryEscape(query))

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "cli-music-player/1.0 (https://github.com/aftaylor2/cli-music-player)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	var result struct {
		Releases []struct {
			ID string `json:"id"`
		} `json:"releases"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Releases) == 0 {
		return "", fmt.Errorf("no releases found")
	}
	return result.Releases[0].ID, nil
}

func downloadCoverArt(mbid string) ([]byte, error) {
	reqURL := fmt.Sprintf("https://coverartarchive.org/release/%s/front-500", mbid)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "cli-music-player/1.0 (https://github.com/aftaylor2/cli-music-player)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	return data, nil
}
