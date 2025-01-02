package database

import (
	"fmt"
	"time"
)

func Get4LatestsTransactions(userId int) ([]WineTransaction, map[int]string, error) {
	query := `
		SELECT 
			wine_transactions.id,
			wine_transactions.wine_id,
			wine_transactions.quantity,
			wine_transactions.type,
			wine_transactions.date,
			wine_wines.name AS wine_name
		FROM wine_transactions
		LEFT JOIN wine_wines ON wine_transactions.wine_id = wine_wines.id
		WHERE wine_wines.account_id = ?
		ORDER BY wine_transactions.date DESC
		LIMIT 4
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var transactions []WineTransaction
	wineIdToName := make(map[int]string)

	for rows.Next() {
		var transaction WineTransaction
		var wineName string
		var dateBytes []byte

		if err := rows.Scan(&transaction.ID, &transaction.WineID, &transaction.Quantity, &transaction.Type, &dateBytes, &wineName); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if len(dateBytes) > 0 {
			transaction.Date, err = time.Parse(wineTransactionTimeFormat, string(dateBytes))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse date: %w", err)
			}
		}

		transactions = append(transactions, transaction)

		if _, exists := wineIdToName[transaction.WineID]; !exists {
			wineIdToName[transaction.WineID] = wineName
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return transactions, wineIdToName, nil
}

func GetAllTransactions(userId int) ([]WineTransaction, error) {
	query := `
		SELECT
			wine_transactions.quantity,
			wine_transactions.type,
			wine_transactions.date
		FROM wine_transactions
		LEFT JOIN wine_wines ON wine_transactions.wine_id = wine_wines.id
		WHERE wine_wines.account_id = ?
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var transactions []WineTransaction

	for rows.Next() {
		var transaction WineTransaction
		var dateBytes []byte

		if err := rows.Scan(&transaction.Quantity, &transaction.Type, &dateBytes); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if len(dateBytes) > 0 {
			transaction.Date, err = time.Parse(wineTransactionTimeFormat, string(dateBytes))
			if err != nil {
				return nil, fmt.Errorf("failed to parse date: %w", err)
			}
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return transactions, nil
}
