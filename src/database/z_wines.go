package database

import (
	"database/sql"
	"fmt"
)

func GetUserWinesCount(userId int) (int, error) {
	var userWinesCount sql.NullInt32
	query := `
		SELECT COUNT(*) 
		FROM wine_wines 
		WHERE account_id = ?
	`

	err := db.QueryRow(query, userId).Scan(&userWinesCount)
	if err != nil {
		return 0, err
	}

	var realUserWinesCount = 0
	if userWinesCount.Valid {
		realUserWinesCount = int(userWinesCount.Int32)
	}

	return realUserWinesCount, nil
}

func GetWinesForDashboard(userId int) (int, int, int, int, error) {
	var totalWines, totalWinesDrankSold, totalBottlesAdded, totalCurrentBottles sql.NullInt32
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
			 WHERE type IN ('added') 
			 AND wine_id IN (SELECT id FROM wine_wines WHERE account_id = ?)) AS total_bottles_added,
			(SELECT SUM(quantity) FROM wine_wines WHERE account_id = ?) AS total_current_bottles`

	err := db.QueryRow(query, userId, userId, userId, userId).
		Scan(&totalWines, &totalWinesDrankSold, &totalBottlesAdded, &totalCurrentBottles)

	if err != nil {
		return 0, 0, 0, 0, err
	}

	var realTotalWines, realTotalWinesDrankSold, realTotalBottlesAdded, realTotalCurrentBottles = 0, 0, 0, 0

	if totalWines.Valid {
		realTotalWines = int(totalWines.Int32)
	}
	if totalWinesDrankSold.Valid {
		realTotalWinesDrankSold = int(totalWinesDrankSold.Int32)
	}
	if totalBottlesAdded.Valid {
		realTotalBottlesAdded = int(totalBottlesAdded.Int32)
	}
	if totalCurrentBottles.Valid {
		realTotalCurrentBottles = int(totalCurrentBottles.Int32)
	}

	return realTotalWines, realTotalWinesDrankSold, realTotalBottlesAdded, realTotalCurrentBottles, nil
}

func GetWinesCountPerRegion(userId int) (map[string]int, error) {
	query := `
		SELECT
			CONCAT(r.name, ' (', r.country, ')') AS region_name,
			SUM(w.quantity) AS total_quantity
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
			SUM(w.quantity) AS total_quantity
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
