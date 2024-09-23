package repository

import (
	"database/sql"
	"music-service/internal/config"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

type Repository struct {
	Authorization
	PlayList
	Song
}

type Authorization interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByUsernameAndPassword(username string, password string) (*models.User, error)
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
}

func NewRepository(db *sql.DB, cfg *config.Config, log *logging.LogrusLogger) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db, cfg, log),
		PlayList:      NewPlayListRepository(db, log),
		Song:          NewSpotifyRepository(db, log),
	}
}
