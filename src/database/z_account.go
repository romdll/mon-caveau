package database

import (
	"database/sql"
	"moncaveau/database/crypt"
)

func CheckIfAccountKeyExists(accountKey string) (bool, int, error) {
	hashedKey := crypt.HashAccountKey(accountKey)
	var accountId int
	err := db.QueryRow("SELECT id FROM accounts WHERE account_key = ?", hashedKey).Scan(&accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, -1, nil
		}
		return false, -1, err
	}
	return true, accountId, nil
}
