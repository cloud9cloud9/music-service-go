package repository

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

type Repository struct {
	Authorization
	PlayList
	Song
	Token
}

type Authorization interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByUsernameAndPassword(username string, password string) (*models.User, error)
}

type Token interface {
	SaveToken(token models.Token) error
	InvalidateToken(userID int) error
	IsTokenValid(token string) (bool, error)
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

func NewRepository(db *sql.DB, log *logging.LogrusLogger) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db, log),
		PlayList:      NewPlayListRepository(db, log),
		Song:          NewSpotifyRepository(db, log),
		Token:         NewTokenRepository(db, log),
	}
}
