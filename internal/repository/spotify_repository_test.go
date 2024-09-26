package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"music-service/internal/models"
	"music-service/pkg/logging"
	"testing"
)

func TestSpotifyRepository_CreateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()

	storage := NewSpotifyRepository(db, logger)

	song := &models.Song{
		ID:          "song123",
		Title:       "Test Song",
		Artist:      "Test Artist",
		Album:       "Test Album",
		AlbumCover:  "cover_url",
		Duration:    200,
		ReleaseDate: "2024-09-26",
		Popularity:  80,
		PreviewURL:  "preview_url",
		ExternalURL: "external_url",
	}

	testCases := []struct {
		name          string
		userId        int
		playlistId    int
		song          *models.Song
		mockSetup     func()
		expectedError error
		expectedID    string
	}{
		{
			name:       "successful song creation",
			userId:     1,
			playlistId: 1,
			song:       song,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec(`^INSERT INTO songs .* ON DUPLICATE KEY UPDATE id=id$`).
					WithArgs(song.ID, song.Title, song.Artist, song.Album, song.AlbumCover, song.Duration, song.ReleaseDate, song.Popularity, song.PreviewURL, song.ExternalURL).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`^INSERT INTO playlist_songs .*`).
					WithArgs(1, song.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
			expectedID:    song.ID,
		},
		{
			name:       "playlist not found",
			userId:     1,
			playlistId: 2,
			song:       song,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: playlistNotFound,
			expectedID:    "",
		},
		{
			name:       "permission denied",
			userId:     1,
			playlistId: 3,
			song:       song,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(2))
			},
			expectedError: permissionDenied,
			expectedID:    "",
		},
		{
			name:       "error inserting song",
			userId:     1,
			playlistId: 1,
			song:       song,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec(`^INSERT INTO songs .* ON DUPLICATE KEY UPDATE id=id$`).
					WithArgs(song.ID, song.Title, song.Artist, song.Album, song.AlbumCover, song.Duration, song.ReleaseDate, song.Popularity, song.PreviewURL, song.ExternalURL).
					WillReturnError(errors.New("failed to insert song"))
			},
			expectedError: errors.New("failed to insert song"),
			expectedID:    "",
		},
		{
			name:       "error adding song to playlist",
			userId:     1,
			playlistId: 1,
			song:       song,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec(`^INSERT INTO songs .* ON DUPLICATE KEY UPDATE id=id$`).
					WithArgs(song.ID, song.Title, song.Artist, song.Album, song.AlbumCover, song.Duration, song.ReleaseDate, song.Popularity, song.PreviewURL, song.ExternalURL).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`^INSERT INTO playlist_songs .*`).
					WithArgs(1, song.ID).
					WillReturnError(errors.New("failed to add song to playlist"))
			},
			expectedError: errors.New("failed to add song to playlist"),
			expectedID:    "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := storage.CreateSong(tt.userId, tt.playlistId, tt.song)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSpotifyRepository_GetAllSongsFromPlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()
	storage := NewSpotifyRepository(db, logger)

	song := &models.Song{
		ID:          "song123",
		Title:       "Test Song",
		Artist:      "Test Artist",
		Album:       "Test Album",
		AlbumCover:  "cover_url",
		Duration:    200,
		ReleaseDate: "2024-09-26",
		Popularity:  80,
		PreviewURL:  "preview_url",
		ExternalURL: "external_url",
	}

	testCases := []struct {
		name          string
		userId        int
		playlistId    int
		mockSetup     func()
		expectedError error
		expectedSongs []*models.Song
	}{
		{
			name:       "successful get all songs from playlist",
			userId:     1,
			playlistId: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT s\.id, s\.title, s\.artist, s\.album, s\.album_cover, s\.duration, 
		                   s\.release_date, s\.popularity, s\.preview_url, s\.external_url
		                   FROM playlist_songs ps
		                   JOIN songs s ON ps.song_id = s.id
		                   JOIN playlists p ON ps.playlist_id = p.id
		                   WHERE p.id = \? AND p.user_id = \?$`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "album_cover", "duration", "release_date", "popularity", "preview_url", "external_url"}).
						AddRow(song.ID, song.Title, song.Artist, song.Album, song.AlbumCover, song.Duration, song.ReleaseDate, song.Popularity, song.PreviewURL, song.ExternalURL))
			},
			expectedError: nil,
			expectedSongs: []*models.Song{song},
		},
		{
			name:       "error getting all songs from playlist",
			userId:     1,
			playlistId: 1,
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT s\.id, s\.title, s\.artist, s\.album, s\.album_cover, s\.duration, 
		                   s\.release_date, s\.popularity, s\.preview_url, s\.external_url
		                   FROM playlist_songs ps
		                   JOIN songs s ON ps.song_id = s.id
		                   JOIN playlists p ON ps.playlist_id = p.id
		                   WHERE p.id = \? AND p.user_id = \?$`).
					WithArgs(1, 1).
					WillReturnError(errors.New("failed to get all songs from playlist"))
			},
			expectedError: errors.New("failed to get all songs from playlist"),
			expectedSongs: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			songs, err := storage.GetAllSongsFromPlaylist(tt.userId, tt.playlistId)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, tt.expectedSongs, songs)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSongs, songs)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSpotifyRepository_DeleteSongFromPlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()
	storage := NewSpotifyRepository(db, logger)

	testCases := []struct {
		name          string
		userId        int
		playlistId    int
		songId        string
		mockSetup     func()
		expectedError error
	}{
		{
			name:       "successful delete song from playlist",
			userId:     1,
			playlistId: 1,
			songId:     "song123",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec(`^DELETE FROM playlist_songs WHERE playlist_id = \? AND song_id = \?$`).
					WithArgs(1, "song123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},

		{
			name:       "error deleting song from playlist",
			userId:     1,
			playlistId: 1,
			songId:     "song123",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

				mock.ExpectExec(`^DELETE FROM playlist_songs WHERE playlist_id = \? AND song_id = \?$`).
					WithArgs(1, "song123").
					WillReturnError(errors.New("failed to delete song from playlist"))
			},
			expectedError: errors.New("failed to delete song from playlist"),
		},

		{
			name:       "user does not own playlist",
			userId:     1,
			playlistId: 1,
			songId:     "song123",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT user_id FROM playlists WHERE id = \?$`).
					WithArgs(1).
					WillReturnError(errors.New("user does not own playlist"))

			},
			expectedError: errors.New("user does not own playlist"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := storage.DeleteSongFromPlaylist(tt.userId, tt.playlistId, tt.songId)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
