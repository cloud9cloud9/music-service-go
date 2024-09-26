package repository

import (
	"database/sql"
	"errors"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

var (
	userNotFound      = errors.New("user not found")
	userAlreadyExists = errors.New("user already exists")
)

type AuthRepository struct {
	storage *sql.DB
	log     *logging.LogrusLogger
}

func NewAuthRepository(
	storage *sql.DB,
	log *logging.LogrusLogger,
) *AuthRepository {
	return &AuthRepository{
		storage: storage,
		log:     log,
	}
}

func (a *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	rows, err := a.storage.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		a.log.Error("REPOSITORY: unsuccessful get user by email: ", err)
		return nil, err
	}

	u := new(models.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			a.log.Error("REPOSITORY: can't scan rows into user: ", err)
			return nil, err
		}
	}

	if u.ID == 0 {
		a.log.Error("REPOSITORY: user not found: ", err)
		return nil, userNotFound
	}

	a.log.Info("REPOSITORY: get user by email: ", u)
	return u, nil
}

func (a *AuthRepository) GetUserByID(id int) (*models.User, error) {
	rows, err := a.storage.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		a.log.Error("REPOSITORY: can't get user by id: ", err)
		return nil, err
	}

	u := new(models.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			a.log.Error("REPOSITORY:  can't scan rows into user: ", err)
			return nil, err
		}
	}

	if u.ID == 0 {
		a.log.Error("REPOSITORY: user not found: ", err)
		return nil, userNotFound
	}

	a.log.Info("REPOSITORY: get user by id: ", u)
	return u, nil
}

func (a *AuthRepository) CreateUser(user *models.User) error {
	_, err := a.storage.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password,
	)
	if err != nil {
		a.log.Error("REPOSITORY: can't create user: ", err)
		return userAlreadyExists
	}

	a.log.Info("REPOSITORY: user created successfully: ", user)
	return nil
}

func (a *AuthRepository) GetUserByUsernameAndPassword(username string, password string) (*models.User, error) {
	rows, err := a.storage.Query(
		"SELECT * FROM users WHERE username = ? AND password = ?",
		username,
		password,
	)
	if err != nil {
		a.log.Error("REPOSITORY: unsuccessful get user by username and password: ", err)
		return nil, err
	}

	user := new(models.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			a.log.Error("REPOSITORY: can't scan rows into user: ", err)
			return nil, err
		}
	}

	if user.ID == 0 {
		a.log.Error("REPOSITORY: user not found: ", err)
		return nil, userNotFound
	}

	a.log.Info("REPOSITORY: get user by username and password: ", username)
	return user, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*models.User, error) {
	var user models.User

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
