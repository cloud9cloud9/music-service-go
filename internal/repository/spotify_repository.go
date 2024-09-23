package repository

import (
	"database/sql"
	"errors"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

type SpotifyRepository struct {
	storage *sql.DB
	log     *logging.LogrusLogger
}

var (
	musicNotFound    = errors.New("no rows in result set")
	playlistNotFound = errors.New("playlist not found")
	permissionDenied = errors.New("user does not own this playlist")
)

func NewSpotifyRepository(
	storage *sql.DB,
	log *logging.LogrusLogger,
) *SpotifyRepository {
	return &SpotifyRepository{
		storage: storage,
		log:     log,
	}
}

func (s *SpotifyRepository) GetAllSongsFromPlaylist(userId, playlistId int) ([]*models.Song, error) {
	var songs []*models.Song
	query := `
		SELECT s.id, s.title, s.artist, s.album, s.album_cover, s.duration, 
		       s.release_date, s.popularity, s.preview_url, s.external_url 
		FROM playlist_songs ps
		JOIN songs s ON ps.song_id = s.id
		JOIN playlists p ON ps.playlist_id = p.id
		WHERE p.id = ? AND p.user_id = ?
	`
	rows, err := s.storage.Query(query, playlistId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Error("REPOSITORY: no songs in playlist")
			return nil, musicNotFound
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		song, err := scanRowsIntoSong(rows)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	s.log.Info("REPOSITORY: get list of tracks from playlist:", len(songs))
	return songs, nil
}

func (s *SpotifyRepository) CreateSong(userId, playlistId int, song *models.Song) (string, error) {
	var playlistOwner int
	query := `SELECT user_id FROM playlists WHERE id = ?`
	err := s.storage.QueryRow(query, playlistId).Scan(&playlistOwner)
	if err != nil {
		s.log.Error("REPOSITORY: get playlist owner:", err)
		if err == sql.ErrNoRows {
			s.log.Error("REPOSITORY: playlist not found:", err)
			return "", playlistNotFound
		}
		return "", err
	}

	if playlistOwner != userId {
		s.log.Error("REPOSITORY: permitting denied:", err)
		return "", permissionDenied
	}

	_, err = s.storage.Exec(`
		INSERT INTO songs (id, title, artist, album, album_cover, duration, release_date, popularity, preview_url, external_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE id=id
	`, song.ID, song.Title, song.Artist, song.Album, song.AlbumCover, song.Duration,
		song.ReleaseDate, song.Popularity, song.PreviewURL, song.ExternalURL)

	if err != nil {
		s.log.Error("REPOSITORY: track not created:", err)
		return "", err
	}

	_, err = s.storage.Exec(`
		INSERT INTO playlist_songs (playlist_id, song_id) 
		VALUES (?, ?)
	`, playlistId, song.ID)

	if err != nil {
		s.log.Error("REPOSITORY: track not added to playlist_songs:", err)
		return "", err
	}

	s.log.Info("REPOSITORY: track created successfully:", song.ID)
	return song.ID, nil
}

func (s *SpotifyRepository) DeleteSongFromPlaylist(userId, playlistId int, songId string) error {
	var playlistOwner int
	query := `SELECT user_id FROM playlists WHERE id = ?`
	err := s.storage.QueryRow(query, playlistId).Scan(&playlistOwner)
	if err != nil {
		s.log.Error("REPOSITORY: get playlist owner:", err)
		if err == sql.ErrNoRows {
			return playlistNotFound
		}
		return err
	}

	if playlistOwner != userId {
		s.log.Error("REPOSITORY: permitting denied:", err)
		return permissionDenied
	}

	_, err = s.storage.Exec(`
		DELETE FROM playlist_songs 
		WHERE playlist_id = ? AND song_id = ?
	`, playlistId, songId)
	if err != nil {
		s.log.Error("REPOSITORY: track not removed from playlist_songs:", err)
		return err
	}

	s.log.Info("REPOSITORY: track removed successfully:", songId)
	return nil
}

func scanRowsIntoSong(rows *sql.Rows) (*models.Song, error) {
	song := new(models.Song)
	err := rows.Scan(
		&song.ID,
		&song.Title,
		&song.Artist,
		&song.Album,
		&song.AlbumCover,
		&song.Duration,
		&song.ReleaseDate,
		&song.Popularity,
		&song.PreviewURL,
		&song.ExternalURL,
	)
	if err != nil {
		return nil, err
	}
	return song, nil
}
