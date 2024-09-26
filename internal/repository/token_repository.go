package repository

import (
	"database/sql"
	"fmt"
	"music-service/internal/models"
	"music-service/pkg/logging"
)

type TokenRepository struct {
	storage *sql.DB
	log     *logging.LogrusLogger
}

func NewTokenRepository(
	storage *sql.DB,
	log *logging.LogrusLogger,
) *TokenRepository {
	return &TokenRepository{
		storage: storage,
		log:     log,
	}
}

func (t *TokenRepository) SaveToken(token models.Token) error {
	var exists bool

	queryCheck := `SELECT EXISTS(SELECT 1 FROM tokens WHERE user_id = ?)`
	err := t.storage.QueryRow(queryCheck, token.UserID).Scan(&exists)
	if err != nil {
		t.log.Error(fmt.Sprintf("Error checking token existence: %s", err))
		return err
	}

	if exists {
		queryUpdate := `
            UPDATE tokens 
            SET token = ?, expires_at = ?, status = 'active'
            WHERE user_id = ?
        `
		_, err := t.storage.Exec(queryUpdate, token.Token, token.ExpiresAt, token.UserID)
		if err != nil {
			t.log.Error(fmt.Sprintf("Error updating token: %s", err))
			return err
		}
		t.log.Info("Token updated successfully")
		return nil
	}

	queryInsert := `
        INSERT INTO tokens (token, expires_at, user_id) 
        VALUES (?, ?, ?)
    `
	_, err = t.storage.Exec(queryInsert, token.Token, token.ExpiresAt, token.UserID)
	if err != nil {
		t.log.Error(fmt.Sprintf("Error saving token: %s", err))
		return err
	}

	t.log.Info("Token saved successfully")
	return nil
}

func (t *TokenRepository) InvalidateToken(userID int) error {
	query := `
        UPDATE tokens 
        SET status = 'inactive' 
        WHERE user_id = ? AND status = 'active'
    `
	_, err := t.storage.Exec(query, userID)
	if err != nil {
		t.log.Error(fmt.Sprintf("Error invalidating tokens for user_id %d: %s", userID, err))
		return err
	}

	t.log.Info(fmt.Sprintf("All active tokens invalidated for user_id %d", userID))
	return nil
}

func (r *TokenRepository) IsTokenValid(token string) (bool, error) {
	var status string
	query := `SELECT status FROM tokens WHERE token = ?`
	err := r.storage.QueryRow(query, token).Scan(&status)
	if err != nil {
		return false, err
	}

	return status == "active", nil
}
