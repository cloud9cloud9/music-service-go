package service

import (
	"github.com/zmb3/spotify"
	"music-service/internal/models"
	"music-service/internal/repository"
)

type Service struct {
	Authorization
	PlayList
	Song
}

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	ParseToken(accessToken string) (int, error)
	CreateToken(username, password string) (string, error)
}

type PlayList interface {
	CreatePlaylist(playlist *models.Playlist) (int64, error)
	GetAllPlaylists(userId int) ([]*models.Playlist, error)
	GetPlaylistById(userId int, playlistId int) (*models.Playlist, error)
	UpdatePlaylistById(userId int, playlist *models.Playlist) error
	DeletePlaylistById(userId int, playlistId int) error
}

type Song interface {
	GetAllSongsFromPlaylist(userId, playlistId int) ([]*models.Song, error)
	CreateSong(userId, playlistId int, song *models.Song) (string, error)
	DeleteSongFromPlaylist(userId, playlistId int, songId string) error
	GetTrackByID(trackID string) (*spotify.FullTrack, error)
}

func NewService(repo *repository.Repository, client *spotify.Client, exp int64, secret string) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, secret, exp),
		PlayList:      NewPlaylistService(repo.PlayList),
		Song:          NewSpotifyService(repo.Song, client),
	}
}
