package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/models"
	"music-service/pkg/logging"
	"testing"
	"time"
)

func TestAuthRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()
	authRepo := NewAuthRepository(db, logger)

	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "test",
		Password:  "test",
		CreatedAt: time.Now(),
	}
	tests := []struct {
		name          string
		email         string
		mockSetup     func()
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:  "successful get user by email",
			email: "test@example.com",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE email = \\?$").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "createdAt"}).
						AddRow(user.ID, user.Password, user.Email, user.Password, user.CreatedAt))
			},
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name:  "error on get user by email",
			email: "test@example.com",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE email = \\?$").
					WithArgs("test@example.com").
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedUser:  nil,
			expectedError: sqlmock.ErrCancelled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := authRepo.GetUserByEmail(tt.email)
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAuthRepository_CreateUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()
	authRepo := NewAuthRepository(db, logger)

	user := &models.User{
		Username: "test",
		Email:    "test@example.com",
		Password: "test",
	}

	tests := []struct {
		name          string
		user          *models.User
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successful user creation",
			user: user,
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Username, user.Email, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "error on user creation",
			user: user,
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test", "test@example.com", "test").
					WillReturnError(errors.New("user already exists"))
			},
			expectedError: errors.New("user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := authRepo.CreateUser(tt.user)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAuthRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()
	authRepo := NewAuthRepository(db, logger)

	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "test",
		Password:  "test",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		id            int
		mockSetup     func()
		expectedUser  *models.User
		expectedError error
	}{
		{
			name: "successful get user by id",
			id:   1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE id = \\?$").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "createdAt"}).
						AddRow(user.ID, user.Username, user.Email, user.Password, user.CreatedAt))
			},
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name: "error on get user by id",
			id:   1,
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE id = \\?$").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := authRepo.GetUserByID(tt.id)
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAuthRepository_GetUserByUsernameAndPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	logger := logging.NewLogger()
	authRepo := NewAuthRepository(db, logger)

	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "test",
		Password:  "test",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		username      string
		password      string
		mockSetup     func()
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:     "successful get user by username and password",
			username: "test",
			password: "test",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE username = \\? AND password = \\?$").
					WithArgs("test", "test").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "createdAt"}).
						AddRow(user.ID, user.Username, user.Email, user.Password, user.CreatedAt))
			},
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name:     "error on get user by username and password",
			username: "test",
			password: "test",
			mockSetup: func() {
				mock.ExpectQuery("^SELECT \\* FROM users WHERE username = \\? AND password = \\?$").
					WithArgs("test", "test").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := authRepo.GetUserByUsernameAndPassword(tt.username, tt.password)
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
