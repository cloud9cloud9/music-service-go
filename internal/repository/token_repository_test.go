package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"music-service/internal/models"
	"music-service/pkg/logging"
	"testing"
	"time"
)

func TestTokenRepository_SaveToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	logger := logging.NewLogger()

	repo := NewTokenRepository(db, logger)

	token := &models.Token{
		Token:     "token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		UserID:    1,
	}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "token exists and is updated successfully",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT EXISTS\(SELECT 1 FROM tokens WHERE user_id = \?\)$`).
					WithArgs(token.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectExec(`^UPDATE tokens SET token = \?, expires_at = \?, status = 'active' WHERE user_id = \?$`).
					WithArgs(token.Token, token.ExpiresAt, token.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "token does not exist and is created successfully",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT EXISTS\(SELECT 1 FROM tokens WHERE user_id = \?\)$`).
					WithArgs(token.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectExec(`^INSERT INTO tokens \(token, expires_at, user_id\) VALUES \(\?, \?, \?\)$`).
					WithArgs(token.Token, token.ExpiresAt, token.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "error during token existence check",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT EXISTS\(SELECT 1 FROM tokens WHERE user_id = \?\)$`).
					WithArgs(token.UserID).
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "error updating token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT EXISTS\(SELECT 1 FROM tokens WHERE user_id = \?\)$`).
					WithArgs(token.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectExec(`^UPDATE tokens SET token = \?, expires_at = \?, status = 'active' WHERE user_id = \?$`).
					WithArgs(token.Token, token.ExpiresAt, token.UserID).
					WillReturnError(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
		{
			name: "error inserting new token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT EXISTS\(SELECT 1 FROM tokens WHERE user_id = \?\)$`).
					WithArgs(token.UserID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectExec(`^INSERT INTO tokens \(token, expires_at, user_id\) VALUES \(\?, \?, \?\)$`).
					WithArgs(token.Token, token.ExpiresAt, token.UserID).
					WillReturnError(errors.New("insert error"))
			},
			expectedError: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.SaveToken(*token)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestTokenRepository_IsTokenValid(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	logger := logging.NewLogger()

	repo := NewTokenRepository(db, logger)

	tests := []struct {
		name          string
		token         string
		mockSetup     func()
		expectedError error
		expectedValid bool
	}{
		{
			name:  "token is valid",
			token: "token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT status FROM tokens WHERE token = \?`).
					WithArgs("token").
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("active"))
			},
			expectedError: nil,
			expectedValid: true,
		},
		{
			name:  "token is not valid",
			token: "token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT status FROM tokens WHERE token = \?`).
					WithArgs("token").
					WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("inactive"))
			},
			expectedError: nil,
			expectedValid: false,
		},
		{
			name:  "error during token existence check",
			token: "token",
			mockSetup: func() {
				mock.ExpectQuery(`^SELECT status FROM tokens WHERE token = \?`).
					WithArgs("token").
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			valid, err := repo.IsTokenValid(tt.token)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedValid, valid)

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestTokenRepository_InvalidateToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	logger := logging.NewLogger()

	repo := NewTokenRepository(db, logger)

	tests := []struct {
		name          string
		userID        int
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "token is invalidated",
			userID: 1,
			mockSetup: func() {
				mock.ExpectExec(`^UPDATE tokens SET status = 'inactive' WHERE user_id = \? AND status = 'active'$`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:   "error during token invalidation",
			userID: 1,
			mockSetup: func() {
				mock.ExpectExec(`^UPDATE tokens SET status = 'inactive' WHERE user_id = \? AND status = 'active'$`).
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.InvalidateToken(tt.userID)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
