package database

import (
	"database/sql"
	"fmt"
)

func GetWinesForDashboard(userId int) (int, int, int, error) {
	var totalWines, totalWinesDrankSold, totalWinesDrankSoldThisMonth sql.NullInt32
	query := `
		SELECT 
			(SELECT COUNT(*) 
			 FROM wine_wines 
			 WHERE account_id = ?) AS total_wines,
			(SELECT SUM(quantity) 
			 FROM wine_transactions 
			 WHERE type IN ('drank', 'sold') 
			 AND wine_id IN (SELECT id FROM wine_wines WHERE account_id = ?)) AS total_wines_drank_sold,
			(SELECT SUM(quantity) 
			 FROM wine_transactions 
			 WHERE type IN ('drank', 'sold') 
			 AND wine_id IN (SELECT id FROM wine_wines WHERE account_id = ?) 
			 AND YEAR(date) = YEAR(CURDATE()) 
			 AND MONTH(date) = MONTH(CURDATE())) AS total_wines_drank_sold_this_month`

	err := db.QueryRow(query, userId, userId, userId).
		Scan(&totalWines, &totalWinesDrankSold, &totalWinesDrankSoldThisMonth)

	if err != nil {
		return 0, 0, 0, err
	}

	var realTotalWines, realTotalWinesDrankSold, realTotalWinesDrankSoldThisMonth = 0, 0, 0

	if totalWines.Valid {
		realTotalWines = int(totalWines.Int32)
	}
	if totalWinesDrankSold.Valid {
		realTotalWinesDrankSold = int(totalWinesDrankSold.Int32)
	}
	if totalWinesDrankSoldThisMonth.Valid {
		realTotalWinesDrankSoldThisMonth = int(totalWinesDrankSoldThisMonth.Int32)
	}

	return realTotalWines, realTotalWinesDrankSold, realTotalWinesDrankSoldThisMonth, nil
}

func GetWinesCountPerRegion(userId int) (map[string]int, error) {
	query := `
		SELECT
			CONCAT(r.name, ' (', r.country, ')') AS region_name,
			SUM(w.qantity) AS total_quantity
		FROM wine_wines w
		JOIN wine_regions r ON w.region_id = r.id
		WHERE w.account_id = ?
		GROUP BY r.id
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	regionWineCount := make(map[string]int)

	for rows.Next() {
		var regionName string
		var totalQuantity int
		if err := rows.Scan(&regionName, &totalQuantity); err != nil {
			return nil, err
		}
		regionWineCount[regionName] = totalQuantity
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return regionWineCount, nil
}

func GetWinesCountPerTypes(userId int) (map[string]int, error) {
	query := `
		SELECT
			wt.name AS wine_type, 
			SUM(w.qantity) AS total_quantity
		FROM wine_wines w
		JOIN wine_types wt ON w.type_id = wt.id
		WHERE w.account_id = ?
		GROUP BY wt.id
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	wineTypeCount := make(map[string]int)

	for rows.Next() {
		var wineType string
		var totalQuantity int
		if err := rows.Scan(&wineType, &totalQuantity); err != nil {
			return nil, err
		}
		wineTypeCount[wineType] = totalQuantity
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wineTypeCount, nil
}

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
		if err := rows.Scan(&transaction.ID, &transaction.WineID, &transaction.Quantity, &transaction.Type, &transaction.Date, &wineName); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
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
