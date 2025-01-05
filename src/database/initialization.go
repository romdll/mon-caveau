package database

import (
	"database/sql"
	"fmt"
	"moncaveau/database/migrations"
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
		logger.Errorw("Failed to open database", "error", err)
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		logger.Errorw("Database connection failed", "error", err)
		return fmt.Errorf("database connection failed: %w", err)
	}

	logger.Infow("Database connected successfully")
	return nil
}

func HealthCheck() bool {
	if err := db.Ping(); err != nil {
		logger.Errorw("Database health check failed", "error", err)
		return false
	}
	return true
}

func CloseDB() {
	if db != nil {
		logger.Infow("Closing database connection")
		err := db.Close()
		if err != nil {
			logger.Errorf("Error when closing the database: %w", err)
		} else {
			logger.Infow("Closed the database connection succesfully")
		}
	}
}

func ApplyMigrations() error {
	return migrations.ApplyMigrations(db)
}

func GetAllTablesAndStructures() ([]migrations.TableInfo, error) {
	return migrations.GetAllTablesAndStructures(db)
}
