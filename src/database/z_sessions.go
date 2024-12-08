package database

import (
	"database/sql"
)

const (
	AuthCookieName = "X-Mon-Caveau-Auth"
)

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
