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
		WITH CountryCheck AS (
			SELECT DISTINCT r.country
			FROM wine_wines w
			JOIN wine_regions r ON w.region_id = r.id
			WHERE w.account_id = ?
		)
		SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM CountryCheck) > 1 THEN CONCAT(r.name, ' (', r.country, ')')
				ELSE r.name
			END AS region_name,
			SUM(w.quantity) AS total_quantity
		FROM wine_wines w
		JOIN wine_regions r ON w.region_id = r.id
		WHERE w.account_id = ?
		GROUP BY r.id
	`

	rows, err := db.Query(query, userId, userId)
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

func GetTop5DomainsPerNumberOfBottles(userId int) (map[string]int, error) {
	query := `
		SELECT wd.name AS domain, wr.name AS region, wr.country, SUM(ww.quantity) AS total_quantity
		FROM wine_wines ww
		JOIN wine_domains wd ON ww.domain_id = wd.id
		JOIN wine_regions wr ON ww.region_id = wr.id
		WHERE ww.account_id = ?
		GROUP BY wd.name, wr.name, wr.country
		ORDER BY total_quantity DESC
		LIMIT 5
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query top 5 wine domains: %w", err)
	}
	defer rows.Close()

	topDomains := make(map[string]int)
	for rows.Next() {
		var domain, region, country string
		var totalQuantity int
		if err := rows.Scan(&domain, &region, &country, &totalQuantity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		key := fmt.Sprintf("%s - %s (%s)", domain, region, country)
		topDomains[key] = totalQuantity
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return topDomains, nil
}

func GetWineDistributionPerVintage(userId int) (map[int]int, error) {
	query := `
		SELECT ww.vintage, SUM(ww.quantity)
		FROM wine_wines ww
		WHERE ww.account_id = ?
		GROUP BY ww.vintage
		ORDER BY ww.vintage
	`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query wine distribution per vintage: %w", err)
	}
	defer rows.Close()

	vintageDistribution := make(map[int]int)
	for rows.Next() {
		var vintage int
		var totalQuantity int
		if err := rows.Scan(&vintage, &totalQuantity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		vintageDistribution[vintage] = totalQuantity
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return vintageDistribution, nil
}

type RegionWithWineTypes struct {
	WineRegion `json:"region"`
	WineTypes  []WineTypeWithCount `json:"wine_types"`
}

type WineTypeWithCount struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetWineTypesDistributionPerRegions(userId int) ([]RegionWithWineTypes, error) {
	query := `
        SELECT 
            wr.id AS region_id, wr.name AS region_name, wr.country,
            wt.id AS type_id, wt.name AS type_name,
            SUM(ww.quantity) AS bottle_count
        FROM wine_wines ww
        JOIN wine_regions wr ON ww.region_id = wr.id
        JOIN wine_types wt ON ww.type_id = wt.id
        WHERE ww.account_id = ?
        GROUP BY wr.id, wr.name, wr.country, wt.id, wt.name
        ORDER BY wr.id, wt.id
    `

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query wine types distribution per regions: %w", err)
	}
	defer rows.Close()

	distribution := make(map[WineRegion][]WineTypeWithCount)
	for rows.Next() {
		var region WineRegion
		var wineType WineTypeWithCount
		if err := rows.Scan(&region.ID, &region.Name, &region.Country, &wineType.ID, &wineType.Name, &wineType.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		distribution[region] = append(distribution[region], wineType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	var result []RegionWithWineTypes
	for region, wineTypes := range distribution {
		result = append(result, RegionWithWineTypes{
			WineRegion: region,
			WineTypes:  wineTypes,
		})
	}

	return result, nil
}

type RegionWithBottleCount struct {
	WineRegion  `json:"region"`
	BottleCount int `json:"bottle_count"`
}

func GetUserUsedRegionsWithBottleCount(userId int) ([]RegionWithBottleCount, error) {
	query := `
        SELECT 
            wr.id AS region_id, wr.name AS region_name, wr.country,
            SUM(ww.quantity) AS total_bottles
        FROM wine_wines ww
        JOIN wine_regions wr ON ww.region_id = wr.id
        WHERE ww.account_id = ?
        GROUP BY wr.id, wr.name, wr.country
        ORDER BY wr.id
    `

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query used regions with bottle count: %w", err)
	}
	defer rows.Close()

	var result []RegionWithBottleCount

	for rows.Next() {
		var region WineRegion
		var bottleCount int
		if err := rows.Scan(&region.ID, &region.Name, &region.Country, &bottleCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, RegionWithBottleCount{
			WineRegion:  region,
			BottleCount: bottleCount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}
