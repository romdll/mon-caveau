package migrations

import (
	"database/sql"
	"fmt"
)

type TableColumn struct {
	ColumnName  string `json:"column_name"`
	DataType    string `json:"data_type"`
	ColumnOrder int    `json:"column_order"`
	IsNullable  string `json:"is_nullable"`
}

type TableInfo struct {
	TableName    string        `json:"table_name"`
	Columns      []TableColumn `json:"columns"`
	RowCount     int           `json:"row_count"`
	CreationTime string        `json:"creation_time"`
	Engine       string        `json:"engine"`
	CharacterSet string        `json:"character_set"`
	Collation    string        `json:"collation"`
}

func GetAllTablesAndStructures(db *sql.DB) ([]TableInfo, error) {
	var tables []TableInfo

	tableQuery := `
		SELECT table_name, table_rows, create_time, engine, table_collation
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
		var (
			tableName       string
			rowCount        int
			creationTimeRaw []byte
			engine          sql.NullString
			tableCollation  sql.NullString
		)

		if err := rows.Scan(&tableName, &rowCount, &creationTimeRaw, &engine, &tableCollation); err != nil {
			logger.Errorf("Failed to scan table metadata: %v", err)
			return nil, fmt.Errorf("failed to scan table metadata: %v", err)
		}

		creationTime := parseSQLDateTime(creationTimeRaw)

		logger.Infof("Processing table: %s", tableName)

		columnQuery := `
			SELECT column_name, data_type, ordinal_position, is_nullable
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

		var columns []TableColumn
		for columnRows.Next() {
			var column TableColumn
			if err := columnRows.Scan(&column.ColumnName, &column.DataType, &column.ColumnOrder, &column.IsNullable); err != nil {
				logger.Errorf("Failed to scan column for table %s: %v", tableName, err)
				return nil, fmt.Errorf("failed to scan column for table %s: %v", tableName, err)
			}
			columns = append(columns, column)
		}

		tables = append(tables, TableInfo{
			TableName:    tableName,
			Columns:      columns,
			RowCount:     rowCount,
			CreationTime: creationTime,
			Engine:       nullStringToString(engine),
			CharacterSet: getCharacterSet(tableCollation.String),
			Collation:    nullStringToString(tableCollation),
		})
		logger.Infof("Added table: %s with %d columns.", tableName, len(columns))
	}

	logger.Infof("Fetched structures for %d tables successfully.", len(tables))
	return tables, nil
}
