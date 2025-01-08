package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
)

var secretKey string

func init() {
	secretKey = os.Getenv("ACCOUNT_KEY_SECRET")
	if secretKey == "" {
		logger.Fatal("ACCOUNT_KEY_SECRET is not set")
	}
}

func HashAccountKey(accountKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(accountKey))
	return hex.EncodeToString(h.Sum(nil))
}

func CustomHashAccountKeysMigration(db *sql.DB) error {
	logger.Info("Starting custom migration to hash account keys")

	rows, err := db.Query("SELECT id, account_key FROM accounts")
	if err != nil {
		logger.Errorf("Failed to fetch account keys: %v", err)
		return fmt.Errorf("failed to fetch account keys: %w", err)
	}
	defer rows.Close()

	tx, err := db.Begin()
	if err != nil {
		logger.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for rows.Next() {
		var id int
		var plaintextKey string

		if err := rows.Scan(&id, &plaintextKey); err != nil {
			tx.Rollback()
			logger.Errorf("Failed to scan account key for account ID %d: %v", id, err)
			return fmt.Errorf("failed to scan account key: %w", err)
		}

		logger.Infof("Hashing account key for account ID %d", id)
		hashedKey := HashAccountKey(plaintextKey)

		_, err = tx.Exec("UPDATE accounts SET account_key = ? WHERE id = ?", hashedKey, id)
		if err != nil {
			tx.Rollback()
			logger.Errorf("Failed to update account key for account ID %d: %v", id, err)
			return fmt.Errorf("failed to update account key for id %d: %w", id, err)
		}

		logger.Infof("Successfully hashed and updated account key for account ID %d", id)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		logger.Errorf("Error iterating over account keys: %v", err)
		return fmt.Errorf("error iterating over account keys: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Successfully completed custom migration to hash account keys")
	return nil
}
