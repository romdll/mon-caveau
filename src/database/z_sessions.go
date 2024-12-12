package database

import (
	"database/sql"
)

const (
	AuthCookieName = "X-Mon-Caveau-Auth"
	IsCookieSecure = false
)

func VerifyIfSessionExistsAndIsValid(sessionToken string) (bool, int, error) {
	var exists bool
	var userID int

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE session_token = ? AND expires_at > NOW()), account_id FROM sessions WHERE session_token = ?", sessionToken, sessionToken).Scan(&exists, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, -1, nil
		}
		return false, -1, err
	}

	return exists, userID, nil
}

func DeleteSessionToken(sessionToken string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	return err
}
