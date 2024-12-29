package database

import (
	"database/sql"
	"moncaveau/utils"
)

var (
	// Database connection
	db *sql.DB

	// Logger
	logger = utils.CreateLogger("database")

	// Time format of wine transactions date
	wineTransactionTimeFormat = "2006-01-02 15:04:05.000000"
)
