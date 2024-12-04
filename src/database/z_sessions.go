package database

import (
	"database/sql"
	"time"
)

const (
	AuthCookieName = "X-Mon-Caveau-Auth"
)

func CreateNewSession(accountKey, sessionToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (account_id, session_token, expires_at)
		SELECT id, ?, ? 
		FROM accounts 
		WHERE account_key = ?
	`

	_, err := db.Exec(query, sessionToken, expiresAt, accountKey)
	if err != nil {
		return err
	}

	return nil
}

func VerifyIfSessionExistsAndIsValid(sessionToken string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE session_token = ? AND expires_at > NOW())", sessionToken).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
