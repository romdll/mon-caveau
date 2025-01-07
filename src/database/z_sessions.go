package database

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	AuthCookieName = "X-Mon-Caveau-Auth"
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

func GetAllUserSessions(userId int) ([]Session, error) {
	query := `
		SELECT
			id,
			account_id,
			SUBSTRING(session_token, 1, 8) AS session_token,
			created_at,
			expires_at,
			last_activity
		FROM sessions
		WHERE account_id = ?
		ORDER BY created_at DESC
	`

	var sessions []Session

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session

		var createdAtRaw, expiresAtRaw, lastActivityRaw []byte
		if err := rows.Scan(&session.ID, &session.AccountID, &session.SessionToken,
			&createdAtRaw,
			&expiresAtRaw,
			&lastActivityRaw); err != nil {
			return nil, err
		}

		session.CreatedAt, err = parseTimeFromBytes(createdAtRaw)
		if err != nil {
			return nil, err
		}

		session.ExpiresAt, err = parseTimeFromBytes(expiresAtRaw)
		if err != nil {
			return nil, err
		}

		session.LastActivity, err = parseTimeFromBytes(lastActivityRaw)
		if err != nil {
			return nil, err
		}

		session.SessionToken += "..."

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func parseTimeFromBytes(data []byte) (time.Time, error) {
	timestampStr := string(data)
	parsedTime, err := time.Parse(time.DateTime, timestampStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %v", err)
	}
	return parsedTime, nil
}
