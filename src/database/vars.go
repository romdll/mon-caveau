package database

import (
	"database/sql"
	"log"
	"moncaveau/utils"
)

var (
	// Database connection
	db *sql.DB

	// Logger
	logger *log.Logger = utils.CreateLogger("database")
)
