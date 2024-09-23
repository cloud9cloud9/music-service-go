package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

var (
	dataNotFound = errors.New("data not found")
)

type PlayListRepository struct {
	storage *sql.DB
	log     *logging.LogrusLogger
}

func NewPlayListRepository(
	storage *sql.DB,
	log *logging.LogrusLogger,
) *PlayListRepository {
	return &PlayListRepository{
		storage: storage,
		log:     log,
	}
}

func (p *PlayListRepository) CreatePlaylist(playlist *models.Playlist) (int64, error) {
	result, err := p.storage.Exec(
		"INSERT INTO playlists (name, user_id) VALUES (?, ?)",
		playlist.Name,
		playlist.UserId,
	)
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful create playlist: ", err)
		return 0, err
	}

	playlistId, err := result.LastInsertId()
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful create playlist! Id is empty: ", err)
		return 0, err
	}

	p.log.Info("REPOSITORY: create playlist: ", playlistId)
	return playlistId, nil
}

func (p *PlayListRepository) GetAllPlaylists(userId int) ([]*models.Playlist, error) {
	rows, err := p.storage.Query(
		"SELECT * FROM playlists WHERE user_id = ?",
		userId,
	)
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful get all playlists: ", err)
		return nil, err
	}
	defer rows.Close()

	var playlists []*models.Playlist

	for rows.Next() {
		playlist, err := scanRowsIntoPlayList(rows)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		p.log.Error("REPOSITORY: unsuccessful get all playlists: ", err)
		return nil, err
	}

	p.log.Info("REPOSITORY:  get list of playlists: ", len(playlists))
	return playlists, nil
}

func (p *PlayListRepository) GetPlaylistById(userId int, playlistId int) (*models.Playlist, error) {
	rows, err := p.storage.Query(
		"SELECT * FROM playlists WHERE user_id = ? AND id = ?",
		userId,
		playlistId,
	)
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful get playlist by id: ", err)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		p.log.Error("REPOSITORY:  unsuccessful get playlist by id: ", err)
		return nil, fmt.Errorf("playlist with id %d not found for user %d", playlistId, userId)
	}

	playlist, err := scanRowsIntoPlayList(rows)
	if err != nil {
		p.log.Error("REPOSITORY:  unsuccessful scan playlist: ", err)
		return nil, err
	}

	p.log.Info("REPOSITORY: get playlist by id: ", playlistId)
	return playlist, nil
}

func (p *PlayListRepository) UpdatePlaylistById(userId int, playlist *models.Playlist) error {
	result, err := p.storage.Exec(
		"UPDATE playlists SET name = ? WHERE user_id = ? AND id = ?",
		playlist.Name,
		userId,
		playlist.ID,
	)
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful update playlist: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful update playlist: ", err)
		return err
	}

	if rowsAffected == 0 {
		p.log.Error("REPOSITORY: unsuccessful update playlist: ", err)
		return dataNotFound
	}

	p.log.Info("REPOSITORY: update playlist, rows affected: ", rowsAffected)
	return nil
}

func (p *PlayListRepository) DeletePlaylistById(userId int, playlistId int) error {
	result, err := p.storage.Exec(
		"DELETE FROM playlists WHERE user_id = ? AND id = ?",
		userId,
		playlistId,
	)
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful delete playlist: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.log.Error("REPOSITORY: unsuccessful delete playlist: ", err)
		return err
	}

	if rowsAffected == 0 {
		p.log.Error("REPOSITORY: unsuccessful delete playlist: ", err)
		return dataNotFound
	}

	p.log.Info("REPOSITORY: delete playlist, rows affected: ", rowsAffected)
	return nil
}

func scanRowsIntoPlayList(rows *sql.Rows) (*models.Playlist, error) {
	playlist := new(models.Playlist)
	err := rows.Scan(
		&playlist.ID,
		&playlist.UserId,
		&playlist.Name,
	)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
