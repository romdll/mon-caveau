package database

import (
	"fmt"
	"reflect"
	"strings"
)

func GetWinesWithPaginationAndSearch(limit, page, userId int, searchQuery string) ([]WineWine, error) {
	if limit <= 0 || page <= 0 {
		logger.Errorw("Invalid limit or page", "limit", limit, "page", page)
		return nil, fmt.Errorf("invalid limit or page: both must be greater than zero")
	}

	logger.Infow("Fetching entities from table", "limit", limit, "page", page, "searchQuery", searchQuery)

	dbName, _, err := entityToMap(new(WineWine))
	if err != nil {
		logger.Errorw("Error in entityToMap", "error", err)
		return nil, err
	}

	queryBase := "SELECT * FROM " + dbName + " WHERE account_id = ?"
	args := []interface{}{userId}

	if searchQuery != "" {
		searchConditions := []string{
			"name LIKE ?",
			"vintage LIKE ?",
			"type_id IN (SELECT id FROM wine_types WHERE name LIKE ?)",
			"region_id IN (SELECT id FROM wine_regions WHERE name LIKE ? OR country LIKE ?)",
			"bottle_size_id IN (SELECT id FROM wine_bottle_sizes WHERE name LIKE ? OR size LIKE ?)",
			"domain_id IN (SELECT id FROM wine_domains WHERE name LIKE ?)",
			"CAST(quantity AS CHAR) LIKE ?",
		}

		var searchClauses []string
		for _, condition := range searchConditions {
			searchClauses = append(searchClauses, condition)

			if strings.Contains(condition, "OR") {
				args = append(args, "%"+searchQuery+"%", "%"+searchQuery+"%")
			} else {
				args = append(args, "%"+searchQuery+"%")
			}
		}

		queryBase += " AND (" + strings.Join(searchClauses, " OR ") + ")"
	}

	queryBase += " LIMIT ? OFFSET ?"
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	logger.Infow("Executing query", "query", queryBase, "args", args)

	rows, err := db.Query(queryBase, args...)
	if err != nil {
		logger.Errorw("Error executing SELECT query", "error", err, "query", queryBase, "args", args)
		return nil, err
	}
	defer rows.Close()

	var entities []WineWine

	for rows.Next() {
		newEntity := new(WineWine)

		var dest []interface{}
		val := reflect.ValueOf(newEntity).Elem()
		numFields := val.NumField()

		for i := 0; i < numFields; i++ {
			field := val.Type().Field(i)

			if field.Name == "DB_NAME" {
				continue
			}

			dest = append(dest, val.Field(i).Addr().Interface())
		}

		err := rows.Scan(dest...)
		if err != nil {
			logger.Errorw("Error scanning row", "error", err)
			return nil, err
		}

		entities = append(entities, *newEntity)
	}

	if err := rows.Err(); err != nil {
		logger.Errorw("Error during rows iteration", "error", err)
		return nil, err
	}

	logger.Infow("Successfully fetched entities", "count", len(entities), "table", dbName)
	return entities, nil
}
