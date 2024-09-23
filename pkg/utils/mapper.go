package utils

import (
	"github.com/zmb3/spotify"
	"music-service/internal/models"
)

const (
	unknownArtist = "Unknown Artist"
	externalURL   = "spotify"
)

func MapTrackToSong(track *spotify.FullTrack) models.Song {
	artist := unknownArtist
	if len(track.Artists) > 0 {
		artist = track.Artists[0].Name
	}

	albumCover := ""
	if len(track.Album.Images) > 0 {
		albumCover = track.Album.Images[0].URL
	}

	return models.Song{
		ID:          string(track.ID),
		Title:       track.Name,
		Artist:      artist,
		Album:       track.Album.Name,
		AlbumCover:  albumCover,
		Duration:    track.Duration / 1000,
		ReleaseDate: track.Album.ReleaseDate,
		Popularity:  track.Popularity,
		PreviewURL:  track.PreviewURL,
		ExternalURL: track.ExternalURLs[externalURL],
	}
}
