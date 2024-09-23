package service

import (
	"music-service/internal/models"
	"music-service/internal/repository"
)

type PlaylistService struct {
	repo repository.PlayList
}

func NewPlaylistService(
	repo repository.PlayList,
) *PlaylistService {
	return &PlaylistService{
		repo: repo,
	}
}

func (p *PlaylistService) CreatePlaylist(playlist *models.Playlist) (int64, error) {
	return p.repo.CreatePlaylist(playlist)
}

func (p *PlaylistService) GetAllPlaylists(userId int) ([]*models.Playlist, error) {
	return p.repo.GetAllPlaylists(userId)
}

func (p *PlaylistService) GetPlaylistById(userId int, playlistId int) (*models.Playlist, error) {
	return p.repo.GetPlaylistById(userId, playlistId)
}

func (p *PlaylistService) UpdatePlaylistById(userId int, playlist *models.Playlist) error {
	return p.repo.UpdatePlaylistById(userId, playlist)
}

func (p *PlaylistService) DeletePlaylistById(userId int, playlistId int) error {
	return p.repo.DeletePlaylistById(userId, playlistId)
}
