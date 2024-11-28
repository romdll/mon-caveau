package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() error {
	dbUser := os.Getenv("MON_CAVEAU_DB_USER")
	dbPassword := os.Getenv("MON_CAVEAU_DB_PASSWORD")
	dbHost := os.Getenv("MON_CAVEAU_DB_HOST")
	dbPort := os.Getenv("MON_CAVEAU_DB_PORT")
	dbName := os.Getenv("MON_CAVEAU_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	logger.Println("Database connected successfully")
	return nil
}

func HealthCheck() bool {
	if err := db.Ping(); err != nil {
		logger.Println("Database health check failed:", err)
		return false
	}
	return true
}

func CloseDB() {
	if db != nil {
		logger.Println("Closing database connection")
		db.Close()
	}
}
