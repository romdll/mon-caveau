package database

import (
	"fmt"
)

type TableColumn struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

type TableInfo struct {
	TableName string        `json:"table_name"`
	Columns   []TableColumn `json:"columns"`
}

func GetAllTablesAndStructures() ([]TableInfo, error) {
	var tables []TableInfo

	tableQuery := `
		SELECT table_name 
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		ORDER BY table_name;
	`

	logger.Info("Starting to fetch tables from the database...")
	rows, err := db.Query(tableQuery)
	if err != nil {
		logger.Errorf("Failed to query tables: %v", err)
		return nil, fmt.Errorf("failed to query tables: %v", err)
	}
	defer rows.Close()
	logger.Info("Successfully queried tables list.")

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			logger.Errorf("Failed to scan table name: %v", err)
			return nil, fmt.Errorf("failed to scan table name: %v", err)
		}
		logger.Infof("Processing table: %s", tableName)

		columnQuery := `
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_name = ? AND table_schema = DATABASE()
			ORDER BY ordinal_position;
		`

		columnRows, err := db.Query(columnQuery, tableName)
		if err != nil {
			logger.Errorf("Failed to query columns for table %s: %v", tableName, err)
			return nil, fmt.Errorf("failed to query columns for table %s: %v", tableName, err)
		}
		defer columnRows.Close()
		logger.Infof("Successfully queried columns for table: %s", tableName)

		var columns []TableColumn
		for columnRows.Next() {
			var column TableColumn
			if err := columnRows.Scan(&column.ColumnName, &column.DataType); err != nil {
				logger.Errorf("Failed to scan column for table %s: %v", tableName, err)
				return nil, fmt.Errorf("failed to scan column for table %s: %v", tableName, err)
			}
			logger.Infof("Found column: %s (%s) in table: %s", column.ColumnName, column.DataType, tableName)
			columns = append(columns, column)
		}

		tables = append(tables, TableInfo{
			TableName: tableName,
			Columns:   columns,
		})
		logger.Infof("Added table: %s with %d columns.", tableName, len(columns))
	}

	logger.Infof("Fetched structures for %d tables successfully.", len(tables))
	return tables, nil
}
