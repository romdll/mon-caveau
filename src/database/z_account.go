package database

import "database/sql"

// TODO check if account as no "force password"
func CheckIfAccountKeyExists(accountKey string) (bool, error) {
	var storedAccountKey string
	err := db.QueryRow("SELECT account_key FROM accounts WHERE account_key = ?", accountKey).Scan(&storedAccountKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
