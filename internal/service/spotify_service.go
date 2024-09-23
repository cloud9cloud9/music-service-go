package service

import (
	"github.com/zmb3/spotify"
	"music-service/internal/models"
	"music-service/internal/repository"
)

type SpotifyService struct {
	repo   repository.Song
	client *spotify.Client
}

func NewSpotifyService(
	repo repository.Song,
	client *spotify.Client,
) *SpotifyService {
	return &SpotifyService{
		repo:   repo,
		client: client,
	}
}

func (s *SpotifyService) GetAllSongsFromPlaylist(userId, playlistId int) ([]*models.Song, error) {
	return s.repo.GetAllSongsFromPlaylist(userId, playlistId)
}

func (s *SpotifyService) CreateSong(userId, playlistId int, song *models.Song) (string, error) {
	return s.repo.CreateSong(userId, playlistId, song)
}

func (s *SpotifyService) DeleteSongFromPlaylist(userId, playlistId int, songId string) error {
	return s.repo.DeleteSongFromPlaylist(userId, playlistId, songId)
}

func (s *SpotifyService) GetTrackByID(trackID string) (*spotify.FullTrack, error) {
	spotifyID := spotify.ID(trackID)
	track, err := s.client.GetTrack(spotifyID)
	if err != nil {
		return nil, err
	}
	return track, nil
}
