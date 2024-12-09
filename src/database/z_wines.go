package database

import "database/sql"

func GetWinesForDashboard(userId int) (int, int, int, error) {
	var totalWines, totalWinesDrankSold, totalWinesDrankSoldThisMonth sql.NullInt32

	err := db.QueryRow(`
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
			 AND MONTH(date) = MONTH(CURDATE())) AS total_wines_drank_sold_this_month`,
		userId, userId, userId).Scan(&totalWines, &totalWinesDrankSold, &totalWinesDrankSoldThisMonth)

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
