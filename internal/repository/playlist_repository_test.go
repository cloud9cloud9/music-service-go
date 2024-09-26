package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"music-service/internal/models"
	"music-service/pkg/logging"
	"testing"
)

func TestPlayListRepository_CreatePlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()
	playlistRepo := NewPlayListRepository(db, logger)

	playlist := &models.Playlist{
		Name:   "My Playlist",
		UserId: 1,
	}

	tests := []struct {
		name          string
		playlist      *models.Playlist
		mockSetup     func()
		expectedID    int64
		expectedError error
	}{
		{
			name:     "successful playlist creation",
			playlist: playlist,
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO playlists").
					WithArgs("My Playlist", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name:     "error on insert",
			playlist: playlist,
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO playlists").
					WithArgs("My Playlist", 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedID:    0,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			playlistId, err := playlistRepo.CreatePlaylist(tt.playlist)

			assert.Equal(t, tt.expectedID, playlistId)
			assert.Equal(t, tt.expectedError, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPlayListRepository_GetAllPlaylists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()

	repo := &PlayListRepository{
		storage: db,
		log:     logger,
	}

	test := []struct {
		name           string
		userId         int
		mockSetup      func()
		expectedError  error
		expectedResult []*models.Playlist
	}{
		{
			name:   "successful playlist get",
			userId: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM playlists WHERE user_id = \\?$").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id"}).
						AddRow(1, 1, "Playlist 1").
						AddRow(2, 1, "Playlist 2"))
			},
			expectedError: nil,
			expectedResult: []*models.Playlist{
				{ID: 1, Name: "Playlist 1", UserId: 1},
				{ID: 2, Name: "Playlist 2", UserId: 1},
			},
		},
		{
			name:   "error getting playlist",
			userId: 1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM playlists WHERE user_id = \\?$").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError:  sql.ErrConnDone,
			expectedResult: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.GetAllPlaylists(tt.userId)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestPlayListRepository_GetPlaylistById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()

	repo := &PlayListRepository{
		storage: db,
		log:     logger,
	}

	playlist := &models.Playlist{
		ID:     1,
		Name:   "Playlist 1",
		UserId: 1,
	}
	tests := []struct {
		name           string
		playlistId     int
		userId         int
		mockSetup      func()
		expectedError  error
		expectedResult *models.Playlist
	}{
		{
			name:       "successful playlist get",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM playlists WHERE user_id = \\? AND id = \\?$").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id"}).
						AddRow(1, 1, "Playlist 1"))
			},
			expectedError:  nil,
			expectedResult: playlist,
		},
		{
			name:       "error getting playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM playlists WHERE user_id = \\? AND id = \\?$").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError:  sql.ErrConnDone,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.GetPlaylistById(tt.playlistId, tt.userId)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestPlayListRepository_UpdatePlaylistById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()

	repo := &PlayListRepository{
		storage: db,
		log:     logger,
	}

	tests := []struct {
		name           string
		playlistId     int
		userId         int
		playlist       *models.Playlist
		mockSetup      func()
		expectedError  error
		expectedResult *models.Playlist
	}{
		{
			name:       "successful playlist update",
			playlistId: 1,
			userId:     1,
			playlist: &models.Playlist{
				Name:   "Playlist 1",
				ID:     1,
				UserId: 1,
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE playlists SET name = \\? WHERE user_id = \\? AND id = \\?$").
					WithArgs("Playlist 1", 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError:  nil,
			expectedResult: &models.Playlist{ID: 1, Name: "Playlist 1", UserId: 1},
		},
		{
			name:       "error updating playlist",
			playlistId: 1,
			userId:     1,
			playlist: &models.Playlist{
				Name: "Playlist 1",
				ID:   1,
			},
			mockSetup: func() {
				mock.ExpectExec("^UPDATE playlists SET name = \\? WHERE user_id = \\? AND id = \\?$").
					WithArgs("Playlist 1", 1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError:  sql.ErrConnDone,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdatePlaylistById(tt.userId, tt.playlist)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, tt.expectedResult)
			} else {
				assert.Equal(t, tt.expectedResult, tt.playlist)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestPlayListRepository_DeletePlaylistById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logging.NewLogger()

	repo := &PlayListRepository{
		storage: db,
		log:     logger,
	}

	tests := []struct {
		name          string
		playlistId    int
		userId        int
		mockSetup     func()
		expectedError error
	}{
		{
			name:       "successful playlist deletion",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				mock.ExpectExec("^DELETE FROM playlists WHERE user_id = \\? AND id = \\?$").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:       "error deleting playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				mock.ExpectExec("^DELETE FROM playlists WHERE user_id = \\? AND id = \\?$").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.DeletePlaylistById(tt.userId, tt.playlistId)

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
