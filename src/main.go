package main

import (
	"log"
	"moncaveau/database"
	"moncaveau/server"
	"moncaveau/utils"
)

var (
	logger *log.Logger = utils.CreateLogger("main")
)

func main() {
	if err := database.InitDB(); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	if err := database.ApplyMigrations(); err != nil {
		logger.Fatalf("Failed to apply migrations: %v", err)
	}

	serverEngine := server.CreateServer()
	utils.RunWithQuitNotification(serverEngine)
}
